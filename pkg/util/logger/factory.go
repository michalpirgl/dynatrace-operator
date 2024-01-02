package logger

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

var Factory = factory{}

type factory struct{}

func (f factory) GetLogger(name string) logr.Logger {
	return ctrl.Log.WithName(name) // Just POC, don't judge, don't want to touch every single file ...
}

func init() {
	ctrl.SetLogger(newLogger()) // Instead of doing this self-made factory thingy, we could just use the controller-runtime, that already has this "storing the logger" mechanism
}
