package command

import "mvc/version"

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	Meta
	Version *version.VersionInfo
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the mvc version"
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Run(_ []string) int {
	c.UI.Output(c.Version.FullVersionNumber(true))
	return 0
}
