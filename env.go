package main

import (
	"go/build"
	"log"
	"strings"
)

var (
	goMobileExe string
)

// "Build flags", used by multiple commands.
var (
	buildA       bool   // -a
	buildI       bool   // -i
	buildN       bool   // -n
	buildV       bool   // -v
	buildX       bool   // -x
	buildO       string // -o
	buildGcflags string // -gcflags
	buildLdflags string // -ldflags
	buildTarget  string // -target
	buildWork    bool   // -work
	customIcon   string // -icon
)

func addBuildFlags(cmd *command) {
	cmd.flag.StringVar(&buildO, "o", "", "")
	cmd.flag.StringVar(&buildGcflags, "gcflags", "", "")
	cmd.flag.StringVar(&buildLdflags, "ldflags", "", "")
	cmd.flag.StringVar(&buildTarget, "target", "desktop", "")

	cmd.flag.BoolVar(&buildA, "a", false, "")
	cmd.flag.BoolVar(&buildI, "i", false, "")
	cmd.flag.Var((*stringsFlag)(&build.Default.BuildTags), "tags", "")
}

func addBuildFlagsNVXWork(cmd *command) {
	cmd.flag.BoolVar(&buildN, "n", false, "")
	cmd.flag.BoolVar(&buildV, "v", false, "")
	cmd.flag.BoolVar(&buildX, "x", false, "")
	cmd.flag.BoolVar(&buildWork, "work", false, "")
}

func addCustomBuildFlags(cmd *command) {
	cmd.flag.StringVar(&customIcon, "icon", "", "")
}

func init() {
	addBuildFlags(cmdBuild)
	addBuildFlagsNVXWork(cmdBuild)
	addCustomBuildFlags(cmdBuild)

	addBuildFlags(cmdInstall)
	addBuildFlagsNVXWork(cmdInstall)
	addCustomBuildFlags(cmdInstall)

	addBuildFlagsNVXWork(cmdInit)

	addBuildFlags(cmdBind)
	addBuildFlagsNVXWork(cmdBind)

	out, err := Output(
		"go", "list", "-f", "{{.Target}}",
		"golang.org/x/mobile/cmd/gomobile")
	if err != nil {
		log.Fatalln(err)
	}
	goMobileExe = strings.Trim(out, " \t\n\r")
}
