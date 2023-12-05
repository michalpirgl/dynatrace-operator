package startup

import (
	"github.com/Dynatrace/dynatrace-operator/pkg/util/logger"
)

type InstallMode string

const (
	AgentInstallerMode InstallMode = "installer"
	AgentCsiMode       InstallMode = "provisioned"

	AgentNoHostTenant                  = "-"
	AgentContainerConfFilenameTemplate = "container_%s.conf"
	AgentInitSecretConfigField         = "config"

	LdPreloadFilename = "ld.so.preload"
	LibAgentProcPath  = "/agent/lib64/liboneagentproc.so"

	AgentBinDirMount      = "/mnt/bin"
	AgentShareDirMount    = "/mnt/share"
	AgentConfigDirMount   = "/mnt/config"
	AgentConfInitDirMount = "/mnt/agent-conf"
)

var (
	log = logger.Factory.GetLogger("injection-startup")
)
