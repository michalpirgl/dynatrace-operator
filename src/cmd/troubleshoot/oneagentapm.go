package troubleshoot

import (
	"errors"

	"github.com/Dynatrace/dynatrace-operator/src/kubeobjects"
)

func checkOneAgentAPM(ctx *TroubleshootContext) error {
	log = newTroubleshootLogger("oneAgentAPM")

	logNewCheckf("checking if OneAgentAPM object exists ...")
	exists, err := kubeobjects.CheckIfOneAgentAPMExists(&ctx.KubeConfig)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("OneAgentAPM object still exists - either delete OneAgentAPM objects or fully install the oneAgent operator")
	}

	logOkf("OneAgentAPM does not exist")
	return nil
}
