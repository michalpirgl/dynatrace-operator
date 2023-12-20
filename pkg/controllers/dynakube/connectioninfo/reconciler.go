package connectioninfo

import (
	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/pkg/api/v1beta1/dynakube"
	dtclient "github.com/Dynatrace/dynatrace-operator/pkg/clients/dynatrace"
	"github.com/Dynatrace/dynatrace-operator/pkg/util/hasher"
	k8ssecret "github.com/Dynatrace/dynatrace-operator/pkg/util/kubeobjects/secret"
	"github.com/Dynatrace/dynatrace-operator/pkg/util/timeprovider"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var NoOneAgentCommunicationHostsError = errors.New("no communication hosts for OneAgent are available")

type Reconciler struct {
	client       client.Client
	apiReader    client.Reader
	dtc          dtclient.Client
	dynakube     *dynatracev1beta1.DynaKube
	scheme       *runtime.Scheme
	timeProvider *timeprovider.Provider
}

func NewReconciler(clt client.Client, apiReader client.Reader, scheme *runtime.Scheme, dynakube *dynatracev1beta1.DynaKube, dtc dtclient.Client) *Reconciler { //nolint:revive // argument-limit doesn't apply to constructors
	return &Reconciler{
		client:       clt,
		apiReader:    apiReader,
		dynakube:     dynakube,
		scheme:       scheme,
		dtc:          dtc,
		timeProvider: timeprovider.New(),
	}
}

func (r Reconciler) updateDynakubeStatus(ctx context.Context) error {
	r.dynakube.Status.UpdatedTimestamp = metav1.Now()
	err := r.client.Status().Update(ctx, r.dynakube)
	if err != nil {
		log.Info("could not update dynakube status", "name", r.dynakube.Name)
		return err
	}
	return nil
}

func (r *Reconciler) Reconcile(ctx context.Context) error {
	oldStatus := r.dynakube.Status.DeepCopy()

	if !r.dynakube.FeatureDisableActivegateRawImage() {
		err := r.reconcileActiveGateConnectionInfo(ctx)
		if err != nil {
			return err
		}
	}

	err := r.reconcileOneAgentConnectionInfo(ctx)
	if err != nil {
		return err
	}

	needStatusUpdate, err := hasher.IsDifferent(oldStatus, r.dynakube.Status)
	if err != nil {
		return errors.WithMessage(err, "failed to compare connection info status hashes")
	} else if needStatusUpdate {
		err = r.updateDynakubeStatus(ctx)
	}

	return err
}

func (r *Reconciler) needsUpdate(ctx context.Context, secretName string, condition *metav1.Condition) (bool, error) {
	query := k8ssecret.NewQuery(ctx, r.client, r.apiReader, log)
	_, err := query.Get(types.NamespacedName{Name: secretName, Namespace: r.dynakube.Namespace})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			log.Info("creating secret, because missing", "secretName", secretName)
			return true, nil
		}
		return false, err
	}
	return isConditionOutdated(r.timeProvider, condition, r.dynakube.FeatureApiRequestThreshold()), nil
}

func (r *Reconciler) reconcileOneAgentConnectionInfo(ctx context.Context) error {
	prevCondition := meta.FindStatusCondition(r.dynakube.Status.Conditions, OneAgentConnectionInfoConditionType)
	needsUpdate, err := r.needsUpdate(ctx, r.dynakube.OneagentTenantSecret(), prevCondition)
	if err != nil {
		setCondition(r.dynakube, OneAgentErrorCondition(err, UnexpectedErrorReason))
		return err
	}
	if !needsUpdate {
		log.Info(dynatracev1beta1.GetCacheValidMessage(
			"OneAgent connection info update",
			r.dynakube.Status.OneAgent.ConnectionInfoStatus.LastRequest,
			r.dynakube.FeatureApiRequestThreshold()))
		return nil
	}

	connectionInfo, err := r.dtc.GetOneAgentConnectionInfo()
	if err != nil {
		err := errors.WithMessage(err, "failed to get OneAgent connection info")
		setCondition(r.dynakube, OneAgentErrorCondition(err, UnexpectedErrorReason))
		return err
	}

	r.updateDynakubeOneAgentStatus(connectionInfo)

	err = r.createTenantTokenSecret(ctx, r.dynakube.OneagentTenantSecret(), connectionInfo.ConnectionInfo)
	if err != nil {
		// TODO: Maybe special condition
		return err
	}

	log.Info("OneAgent connection info updated")

	if len(connectionInfo.Endpoints) == 0 {
		log.Info("tenant has no endpoints", "tenant", connectionInfo.TenantUUID)
	}

	if len(connectionInfo.CommunicationHosts) == 0 {
		log.Info("no OneAgent communication hosts received, tenant API requests not yet throttled")
		err := NoOneAgentCommunicationHostsError
		setCondition(r.dynakube, OneAgentErrorCondition(err, NoCommunicationHostsErrorReason))
		return err
	}

	log.Info("received OneAgent communication hosts", "communication hosts", connectionInfo.CommunicationHosts, "tenant", connectionInfo.TenantUUID)
	setCondition(r.dynakube, OneAgentReadyCondition())
	return nil
}

func (r *Reconciler) updateDynakubeOneAgentStatus(connectionInfo dtclient.OneAgentConnectionInfo) {
	r.dynakube.Status.OneAgent.ConnectionInfoStatus.TenantUUID = connectionInfo.TenantUUID
	r.dynakube.Status.OneAgent.ConnectionInfoStatus.Endpoints = connectionInfo.Endpoints
	copyCommunicationHosts(&r.dynakube.Status.OneAgent.ConnectionInfoStatus, connectionInfo.CommunicationHosts)
}

func copyCommunicationHosts(dest *dynatracev1beta1.OneAgentConnectionInfoStatus, src []dtclient.CommunicationHost) {
	dest.CommunicationHosts = make([]dynatracev1beta1.CommunicationHostStatus, 0, len(src))
	for _, host := range src {
		dest.CommunicationHosts = append(dest.CommunicationHosts, dynatracev1beta1.CommunicationHostStatus{
			Protocol: host.Protocol,
			Host:     host.Host,
			Port:     host.Port,
		})
	}
}

func (r *Reconciler) reconcileActiveGateConnectionInfo(ctx context.Context) error {
	prevCondition := meta.FindStatusCondition(r.dynakube.Status.Conditions, ActiveGateConnectionInfoConditionType)
	needsUpdate, err := r.needsUpdate(ctx, r.dynakube.ActivegateTenantSecret(), prevCondition)
	if err != nil {
		setCondition(r.dynakube, ActiveGateErrorCondition(err, UnexpectedErrorReason))
		return err

	}
	if !needsUpdate {
		log.Info(dynatracev1beta1.GetCacheValidMessage(
			"activegate connection info update",
			r.dynakube.Status.ActiveGate.ConnectionInfoStatus.LastRequest,
			r.dynakube.FeatureApiRequestThreshold()))
		return nil
	}

	connectionInfo, err := r.dtc.GetActiveGateConnectionInfo()
	if err != nil {
		log.Info("failed to get activegate connection info")
		setCondition(r.dynakube, ActiveGateErrorCondition(err, UnexpectedErrorReason))
		return err
	}

	r.updateDynakubeActiveGateStatus(connectionInfo)

	err = r.createTenantTokenSecret(ctx, r.dynakube.ActivegateTenantSecret(), connectionInfo.ConnectionInfo)
	if err != nil {
		// TODO: Maybe special condition for this error?
		return err
	}

	log.Info("activegate connection info updated")
	setCondition(r.dynakube, ActiveGateReadyCondition())
	return nil
}

func (r *Reconciler) updateDynakubeActiveGateStatus(connectionInfo dtclient.ActiveGateConnectionInfo) {
	r.dynakube.Status.ActiveGate.ConnectionInfoStatus.TenantUUID = connectionInfo.TenantUUID
	r.dynakube.Status.ActiveGate.ConnectionInfoStatus.Endpoints = connectionInfo.Endpoints
}

func (r *Reconciler) createTenantTokenSecret(ctx context.Context, secretName string, connectionInfo dtclient.ConnectionInfo) error {
	secretData := extractSensitiveData(connectionInfo)
	secret, err := k8ssecret.Create(r.scheme, r.dynakube,
		k8ssecret.NewNameModifier(secretName),
		k8ssecret.NewNamespaceModifier(r.dynakube.Namespace),
		k8ssecret.NewDataModifier(secretData))
	if err != nil {
		return errors.WithStack(err)
	}

	query := k8ssecret.NewQuery(ctx, r.client, r.apiReader, log)
	err = query.CreateOrUpdate(*secret)
	if err != nil {
		log.Info("could not create or update secret for connection info", "name", secret.Name)
		return err
	}
	return nil
}

func extractSensitiveData(connectionInfo dtclient.ConnectionInfo) map[string][]byte {
	data := map[string][]byte{
		TenantTokenName: []byte(connectionInfo.TenantToken),
	}
	return data
}
