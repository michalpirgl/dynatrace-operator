package checks

import (
	"errors"
	"github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot"

	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects"
)

func checkOneAgentAPM(ctx *troubleshoot.TroubleshootContext) error {
	log = troubleshoot.newTroubleshootLogger("oneAgentAPM")

	troubleshoot.LogNewCheckf("checking if OneAgentAPM object exists ...")
	exists, err := kubeobjects.CheckIfOneAgentAPMExists(&ctx.KubeConfig)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("OneAgentAPM object still exists - either delete OneAgentAPM objects or fully install the oneAgent operator")
	}

	troubleshoot.LogOkf("OneAgentAPM does not exist")
	return nil
}
