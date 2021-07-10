package flatyaml

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	version = "v0.1-dev"

	// overwrite version when tagVersion exists
	tagVersion = ""

	// gitCommit is the git sha1
	gitCommit = ""

	// gitTreeState is the state of the git tree {dirty or clean}
	gitTreeState = ""
)

type BuildInfo struct {
	Version      string
	GitCommit    string
	GitTreeState string
	GoVersion    string
}

func GetVersion() string {
	info := BuildInfo{
		Version:      version,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		GoVersion:    runtime.Version(),
	}

	if tagVersion != "" {
		info.Version = tagVersion
	}

	data, _ := json.Marshal(info)
	return fmt.Sprintf("version.BuildInfo%s", string(data))
}
