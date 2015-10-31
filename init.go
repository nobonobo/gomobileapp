package main

var cmdInit = &command{
	run:   runInit,
	Name:  "init",
	Usage: "[-u]",
	Short: "install android compiler toolchain",
	Long: `
Init installs the Android C++ compiler toolchain and builds copies
of the Go standard library for mobile devices.

When first run, it downloads part of the Android NDK.
The toolchain is installed in $GOPATH/pkg/gomobile.

The -u option forces download and installation of the new toolchain
even when the toolchain exists.
`,
}

var initU bool // -u

func init() {
	cmdInit.flag.BoolVar(&initU, "u", false, "force toolchain download")
}

func runInit(cmd *command) error {
	return Cmd(CmdContext{}, goMobileExe, append([]string{"init"}, cmd.flag.Args()...)...)
}
