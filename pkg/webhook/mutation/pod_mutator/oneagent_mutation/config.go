package oneagent_mutation

import (
	"github.com/Dynatrace/dynatrace-operator/pkg/util/logger"
)

var (
	log = logger.Factory.GetLogger("oneagent-pod-mutation")
)

const (
	preloadEnv     = "LD_PRELOAD"
	networkZoneEnv = "DT_NETWORK_ZONE"

	// TODO
	proxyEnv = "DT_PROXY"

	// TODO
	dynatraceMetadataEnv = "DT_DEPLOYMENT_METADATA"

	// TODO
	releaseVersionEnv = "DT_RELEASE_VERSION"

	// TODO
	releaseProductEnv = "DT_RELEASE_PRODUCT"

	// TODO
	releaseStageEnv = "DT_RELEASE_STAGE"

	// TODO
	releaseBuildVersionEnv = "DT_RELEASE_BUILD_VERSION"

	OneAgentBinVolumeName     = "oneagent-bin"
	oneAgentShareVolumeName   = "oneagent-share"
	injectionConfigVolumeName = "injection-config"

	oneAgentCustomKeysPath = "/var/lib/dynatrace/oneagent/agent/customkeys"

	// TODO
	customCertFileName = "custom.pem"

	preloadPath       = "/etc/ld.so.preload"
	containerConfPath = "/var/lib/dynatrace/oneagent/agent/config/container.conf"

	// readonly CSI
	oneagentConfVolumeName = "oneagent-agent-conf"
	OneAgentConfMountPath  = "/opt/dynatrace/oneagent-paas/agent/conf"

	oneagentDataStorageVolumeName = "oneagent-data-storage"
	oneagentDataStorageMountPath  = "/opt/dynatrace/oneagent-paas/datastorage"

	oneagentLogVolumeName = "oneagent-log"

	// TODO
	oneagentLogMountPath = "/opt/dynatrace/oneagent-paas/log"
)
