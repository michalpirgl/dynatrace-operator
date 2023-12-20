package connectioninfo

import (
	"time"

	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/pkg/api/v1beta1/dynakube"
	"github.com/Dynatrace/dynatrace-operator/pkg/util/timeprovider"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	OneAgentConnectionInfoConditionType   = "OneAgentConnectionInfo"
	ActiveGateConnectionInfoConditionType = "ActiveGateConnectionInfo"

	UpToDateConnectionInfoReason    = "UpToDate"
	NoCommunicationHostsErrorReason = "NoCommunicationHosts"
	UnexpectedErrorReason           = "UnexpectedError" // TODO: more specific error reason could exist
)

func ActiveGateReadyCondition() metav1.Condition {
	return metav1.Condition{
		Type:   ActiveGateConnectionInfoConditionType,
		Status: metav1.ConditionTrue,
		Reason: UpToDateConnectionInfoReason,
	}
}

func ActiveGateErrorCondition(err error, reason string) metav1.Condition {
	return metav1.Condition{
		Type:    ActiveGateConnectionInfoConditionType,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: err.Error(),
	}
}

func OneAgentReadyCondition() metav1.Condition {
	return metav1.Condition{
		Type:   OneAgentConnectionInfoConditionType,
		Status: metav1.ConditionTrue,
		Reason: UpToDateConnectionInfoReason,
	}
}

func OneAgentErrorCondition(err error, reason string) metav1.Condition {
	return metav1.Condition{
		Type:    OneAgentConnectionInfoConditionType,
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: err.Error(),
	}
}

func setCondition(dynakube *dynatracev1beta1.DynaKube, newCondition metav1.Condition) {
	newCondition.LastTransitionTime = metav1.Now() // TODO: use timeprovider ?
	meta.SetStatusCondition(&dynakube.Status.Conditions, newCondition)
}

// TODO: maybe a general helper?
func isConditionOutdated(timeProvider *timeprovider.Provider, condition *metav1.Condition, timeout time.Duration) bool {
	if condition == nil && condition.Status == metav1.ConditionFalse { // If not set or failed previously, we need to update
		return true
	}
	return timeProvider.IsOutdated(&condition.LastTransitionTime, timeout)
}
