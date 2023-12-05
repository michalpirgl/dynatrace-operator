package startup

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnv(t *testing.T) {
	t.Run(`create new env for oneagent and data-ingest injection`, func(t *testing.T) {
		resetEnv := prepCombinedTestEnv(t)

		env, err := newEnv()
		resetEnv()

		require.NoError(t, err)
		require.NotNil(t, env)

		assert.Equal(t, AgentCsiMode, env.Mode)
		assert.Equal(t, FailFailurePolicy, env.FailurePolicy)
		assert.NotEmpty(t, env.InstallerFlavor)
		assert.NotEmpty(t, env.InstallerTech)
		assert.NotEmpty(t, env.InstallPath)
		assert.NotEmpty(t, env.InstallVersion)
		assert.Len(t, env.Containers, 5)

		assert.NotEmpty(t, env.K8NodeName)
		assert.NotEmpty(t, env.K8PodName)
		assert.NotEmpty(t, env.K8PodUID)
		assert.NotEmpty(t, env.K8BasePodName)
		assert.NotEmpty(t, env.K8Namespace)
		assert.NotEmpty(t, env.K8ClusterID)

		assert.NotEmpty(t, env.WorkloadKind)
		assert.NotEmpty(t, env.WorkloadName)

		assert.True(t, env.OneAgentInjected)
		assert.True(t, env.DataIngestInjected)
		assert.True(t, env.IsReadOnlyCSI)
	})
	t.Run(`create new env for only data-ingest injection`, func(t *testing.T) {
		resetEnv := prepDataIngestTestEnv(t, false)

		env, err := newEnv()
		resetEnv()

		require.NoError(t, err)
		require.NotNil(t, env)

		assert.Empty(t, env.Mode)
		assert.Equal(t, FailFailurePolicy, env.FailurePolicy)
		assert.NotEmpty(t, env.InstallerFlavor) // set to what is defined in arch.Flavor
		assert.Empty(t, env.InstallerTech)
		assert.Empty(t, env.InstallVersion)
		assert.Empty(t, env.InstallPath)
		assert.Empty(t, env.Containers)

		assert.Empty(t, env.K8NodeName)
		assert.Empty(t, env.K8BasePodName)
		assert.NotEmpty(t, env.K8PodName)
		assert.NotEmpty(t, env.K8PodUID)
		assert.NotEmpty(t, env.K8Namespace)

		assert.NotEmpty(t, env.K8ClusterID)
		assert.NotEmpty(t, env.WorkloadKind)
		assert.NotEmpty(t, env.WorkloadName)

		assert.False(t, env.OneAgentInjected)
		assert.True(t, env.DataIngestInjected)
	})
	t.Run(`create new env for only data-ingest injection with unknown owner workload`, func(t *testing.T) {
		resetEnv := prepDataIngestTestEnv(t, true)

		env, err := newEnv()
		resetEnv()

		require.NoError(t, err)
		require.NotNil(t, env)

		assert.NotEmpty(t, env.K8ClusterID)
		assert.Empty(t, env.WorkloadKind)
		assert.Empty(t, env.WorkloadName)

		assert.False(t, env.OneAgentInjected)
		assert.True(t, env.DataIngestInjected)
	})
	t.Run(`create new env for only oneagent`, func(t *testing.T) {
		resetEnv := prepOneAgentTestEnv(t)

		env, err := newEnv()
		resetEnv()

		require.NoError(t, err)
		require.NotNil(t, env)

		assert.Equal(t, AgentCsiMode, env.Mode)
		assert.Equal(t, FailFailurePolicy, env.FailurePolicy)
		assert.NotEmpty(t, env.InstallerFlavor)
		assert.NotEmpty(t, env.InstallerTech)
		assert.NotEmpty(t, env.InstallVersion)
		assert.NotEmpty(t, env.InstallPath)
		assert.Len(t, env.Containers, 5)

		assert.NotEmpty(t, env.K8NodeName)
		assert.NotEmpty(t, env.K8PodName)
		assert.NotEmpty(t, env.K8PodUID)
		assert.NotEmpty(t, env.K8BasePodName)
		assert.NotEmpty(t, env.K8Namespace)

		assert.NotEmpty(t, env.K8ClusterID)
		assert.Empty(t, env.WorkloadKind)
		assert.Empty(t, env.WorkloadName)

		assert.True(t, env.OneAgentInjected)
		assert.False(t, env.DataIngestInjected)
		assert.True(t, env.IsReadOnlyCSI)
	})
}

func TestFailurePolicyModes(t *testing.T) {
	modes := map[string]string{
		FailFailurePolicy:   FailFailurePolicy,
		SilentFailurePolicy: SilentFailurePolicy,
		ForceFailurePolicy:  ForceFailurePolicy,
		"Fail":              SilentFailurePolicy,
		"other":             SilentFailurePolicy,
	}
	for configuredMode, expectedMode := range modes {
		t.Run(`injection failure policy: `+configuredMode, func(t *testing.T) {
			resetEnv := prepDataIngestTestEnv(t, true)

			err := os.Setenv(InjectionFailurePolicyEnv, configuredMode)
			require.NoError(t, err)

			env, err := newEnv()
			resetEnv()

			require.NoError(t, err)
			require.NotNil(t, env)

			assert.Equal(t, expectedMode, env.FailurePolicy)
		})
	}
}

func prepCombinedTestEnv(t *testing.T) func() {
	resetDataIngestEnvs := prepDataIngestTestEnv(t, false)
	resetOneAgentEnvs := prepOneAgentTestEnv(t)
	return func() {
		resetDataIngestEnvs()
		resetOneAgentEnvs()
	}
}

func prepOneAgentTestEnv(t *testing.T) func() {
	envs := []string{
		AgentInstallerFlavorEnv,
		AgentInstallerTechEnv,
		AgentInstallerVersionEnv,
		K8sNodeNameEnv,
		K8sPodNameEnv,
		K8sPodUIDEnv,
		K8sBasePodNameEnv,
		K8sNamespaceEnv,
		AgentInstallPathEnv,
		K8sClusterIDEnv,
	}
	for i := 1; i <= 5; i++ {
		envs = append(envs, fmt.Sprintf(AgentContainerNameEnvTemplate, i))
		envs = append(envs, fmt.Sprintf(AgentContainerImageEnvTemplate, i))
	}
	for _, envvar := range envs {
		err := os.Setenv(envvar, fmt.Sprintf("TEST_%s", envvar))
		require.NoError(t, err)
	}

	// Int env
	envs = append(envs, AgentContainerCountEnv)
	err := os.Setenv(AgentContainerCountEnv, "5")
	require.NoError(t, err)

	// Mode Env
	envs = append(envs, InjectionFailurePolicyEnv)
	err = os.Setenv(InjectionFailurePolicyEnv, "fail")
	require.NoError(t, err)
	envs = append(envs, AgentInstallModeEnv)
	err = os.Setenv(AgentInstallModeEnv, string(AgentCsiMode))
	require.NoError(t, err)

	// Bool envs
	err = os.Setenv(AgentInjectedEnv, trueStatement)
	require.NoError(t, err)
	envs = append(envs, AgentInjectedEnv)
	err = os.Setenv(AgentReadonlyCSIEnv, trueStatement)
	require.NoError(t, err)
	envs = append(envs, AgentReadonlyCSIEnv)

	return resetTestEnv(envs)
}

func prepDataIngestTestEnv(t *testing.T, isUnknownWorkload bool) func() {
	envs := []string{
		EnrichmentWorkloadKindEnv,
		EnrichmentWorkloadNameEnv,
		K8sClusterIDEnv,
		K8sPodNameEnv,
		K8sPodUIDEnv,
		K8sNamespaceEnv,
	}
	for _, envvar := range envs {
		if isUnknownWorkload &&
			(envvar == EnrichmentWorkloadKindEnv || envvar == EnrichmentWorkloadNameEnv) {
			err := os.Setenv(envvar, "UNKNOWN")
			require.NoError(t, err)
		} else {
			err := os.Setenv(envvar, fmt.Sprintf("TEST_%s", envvar))
			require.NoError(t, err)
		}
	}

	// Mode Env
	envs = append(envs, InjectionFailurePolicyEnv)
	err := os.Setenv(InjectionFailurePolicyEnv, "fail")
	require.NoError(t, err)

	// Bool envs
	err = os.Setenv(EnrichmentInjectedEnv, "true")
	require.NoError(t, err)
	envs = append(envs, EnrichmentInjectedEnv)

	return resetTestEnv(envs)
}

func resetTestEnv(envs []string) func() {
	return func() {
		for _, envvar := range envs {
			_ = os.Unsetenv(envvar)
		}
	}
}
