package version

import (
	"bytes"
	"fmt"
)

var (
	GitCommit        string // This will be filled in by the compiler.
	GitDescribe      string // This will be filled in by the compiler.
	Prerelease       = "dev"
	RevisionMetadata = "" // further describing the build type.
)

// VersionInfo
type VersionInfo struct {
	Revision         string
	Prerelease       string
	RevisionMetadata string
}

func GetVersion() *VersionInfo {
	rel := Prerelease
	md := RevisionMetadata

	if GitDescribe == "" && rel == "" && Prerelease != "" {
		rel = "dev"
	}

	return &VersionInfo{
		Revision:         GitCommit,
		Prerelease:       rel,
		RevisionMetadata: md,
	}
}

func (c *VersionInfo) VersionNumber() string {
	version := fmt.Sprintf("%s", c.Revision)[0:7]

	if c.Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.Prerelease)
	}

	if c.RevisionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, c.RevisionMetadata)
	}

	return version
}

func (c *VersionInfo) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "workspace v%s", c.Revision[0:7])

	if c.Prerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.Prerelease)
	}

	if c.RevisionMetadata != "" {
		fmt.Fprintf(&versionString, "+%s", c.RevisionMetadata)
	}

	return versionString.String()
}
