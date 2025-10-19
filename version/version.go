package version

import (
	"fmt"
	"runtime"
)

var GitCommit string

const Version = "0.1.0"

var BuildDate = ""

var GoVersion = runtime.Version()

var OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)

var AppName = "myapp"

var Branch = ""

type V struct {
	GitCommit string `json:"git_commit"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	OSArch    string `json:"os_arch"`
	AppName   string `json:"app"`
	Branch    string `json:"branch"`
}

func Get() V {
	return V{
		GitCommit: GitCommit,
		Version:   Version,
		BuildDate: BuildDate,
		GoVersion: GoVersion,
		OSArch:    OsArch,
		AppName:   AppName,
		Branch:    Branch,
	}
}
