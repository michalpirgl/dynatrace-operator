package troubleshoot

import "github.com/Dynatrace/dynatrace-operator/src/api/v1beta1"

type Component string

const (
	ComponentOneAgent    Component = "OneAgent"
	ComponentCodeModules Component = "OneAgentCodeModules"
	ComponentActiveGate  Component = "ActiveGate"

	customImagePostfix = " (custom image)"
)

func (c Component) String() string {
	return string(c)
}

func (c Component) Name(isCustomImage bool) string {
	if isCustomImage {
		return c.String() + customImagePostfix
	}
	return c.String()
}

func (c Component) SkipImageCheck(image string) bool {
	return image == "" && c != ComponentCodeModules
}

func (c Component) GetImage(dynaKube *v1beta1.DynaKube) (string, bool) {
	if dynaKube == nil {
		return "", false
	}

	switch c {
	case ComponentOneAgent:
		return dynaKube.OneAgentImage(), dynaKube.CustomOneAgentImage() != ""
	case ComponentCodeModules:
		return dynaKube.CodeModulesImage(), true
	case ComponentActiveGate:
		activeGateImage := dynaKube.ActiveGateImage()
		return activeGateImage, activeGateImage != ""
	}
	return "", false
}
