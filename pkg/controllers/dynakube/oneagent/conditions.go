package oneagent

import (
	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/pkg/api/v1beta1/dynakube"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DaemonSetDeploymentConditionType = "OneAgentDaemonSetDeployment"

	DaemonSetCreatedReason = "DaemonSetCreated" // TODO: more specific reason could exist, like DaemonSetReady
	UnexpectedErrorReason  = "UnexpectedError"  // TODO: more specific error reason could exist
)

func DaemonSetDeploymentCreatedCondition() metav1.Condition {
	return metav1.Condition{
		Type:   DaemonSetDeploymentConditionType,
		Status: metav1.ConditionTrue,
		Reason: DaemonSetCreatedReason,
	}
}

func DaemonSetDeploymentErrorCondition(err error) metav1.Condition {
	return metav1.Condition{
		Type:    DaemonSetDeploymentConditionType,
		Status:  metav1.ConditionFalse,
		Reason:  UnexpectedErrorReason,
		Message: err.Error(),
	}
}

func setCondition(dynakube *dynatracev1beta1.DynaKube, newCondition metav1.Condition) {
	newCondition.LastTransitionTime = metav1.Now() // TODO: use timeprovider ?
	meta.SetStatusCondition(&dynakube.Status.Conditions, newCondition)
}
