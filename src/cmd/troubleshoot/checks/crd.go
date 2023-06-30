package checks

import (
	dynatracev1beta1 "github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"
	"github.com/Dynatrace/dynatrace-operator/src/cmd/troubleshoot"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func checkCRD(troubleshootCtx *troubleshoot.TroubleshootContext) error {
	log = troubleshoot.newTroubleshootLogger("crd")

	troubleshoot.LogNewCheckf("checking if CRD for Dynakube exists ...")

	dynakubeList := &dynatracev1beta1.DynaKubeList{}
	err := troubleshootCtx.ApiReader.List(troubleshootCtx.Context, dynakubeList, &client.ListOptions{Namespace: troubleshootCtx.Namespace})

	if err != nil {
		return determineDynakubeError(err)
	}

	troubleshoot.LogOkf("CRD for Dynakube exists")
	return nil
}

func determineDynakubeError(err error) error {
	if runtime.IsNotRegisteredError(err) {
		err = errors.Wrap(err, "CRD for Dynakube missing")
	} else {
		err = errors.Wrap(err, "could not list Dynakube")
	}
	return err
}
