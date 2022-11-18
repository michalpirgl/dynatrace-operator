package server

import (
	"github.com/Dynatrace/dynatrace-operator/src/cmd/manager/mocks"
	"testing"

	mocks "github.com/Dynatrace/dynatrace-operator/mocks/cmd/manager"
	"github.com/Dynatrace/dynatrace-operator/src/cmd/config"
	dtcsi "github.com/Dynatrace/dynatrace-operator/src/controllers/csi"
	"github.com/Dynatrace/dynatrace-operator/src/logger"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/rest"
)

func TestCsiCommand(t *testing.T) {
	configProvider := &config.MockProvider{}
	configProvider.On("GetConfig").Return(&rest.Config{}, nil)

	cmdMgr := &mocks.Manager{}

	managerProvider := &mocks.Provider{}
	managerProvider.On("CreateManager", mock.Anything, mock.Anything).Return(cmdMgr, nil)

	memFs := afero.NewMemMapFs()
	builder := NewCsiServerCommandBuilder().
		SetConfigProvider(configProvider).
		setManagerProvider(managerProvider).
		SetNamespace("test-namespace").
		setFilesystem(memFs)
	command := builder.Build()
	commandFn := builder.buildRun()

	err := commandFn(command, make([]string, 0))

	// sqlite library does not use afero fs, so it throws an error because path does not exist
	assert.Error(t, err)
	configProvider.AssertCalled(t, "GetConfig")
	managerProvider.AssertCalled(t, "CreateManager", "test-namespace", &rest.Config{})

	exists, err := afero.Exists(memFs, dtcsi.DataPath)
	assert.True(t, exists)
	assert.NoError(t, err)

	// Logging a newline because otherwise `go test` doesn't recognize the result
	logger.Factory.GetLogger("csi command").Info("")
}
