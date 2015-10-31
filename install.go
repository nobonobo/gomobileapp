// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

var cmdInstall = &command{
	run:   runInstall,
	Name:  "install",
	Usage: "[-target android|ios] [-icon icon.png] [-o output] [build flags] [package]",
	Short: "compile android APK and install on device",
	Long: `
Install compiles and installs the app named by the import path on the
attached mobile device.

Only -target android is supported. The 'adb' tool must be on the PATH.

The build flags -a, -i, -n, -x, -gcflags, -ldflags, -tags, and -work are
shared with the build command.
For documentation, see 'go help build'.
`,
}

func runInstall(cmd *command) (err error) {
	var pkg *build.Package
	cwd, _ := os.Getwd()
	switch len(cmd.flag.Args()) {
	case 0:
		pkg, err = build.Default.ImportDir(cwd, build.ImportComment)
	case 1:
		pkg, err = build.Default.Import(cmd.flag.Args()[0], cwd, build.ImportComment)
	default:
		cmd.usage()
		os.Exit(1)
	}
	if err != nil {
		return err
	}
	if pkg.Name != "main" && buildO != "" {
		return fmt.Errorf("cannot set -o when building non-main package")
	}

	switch buildTarget {
	case "android":
		if err := runBuild(cmd); err != nil {
			return err
		}
		name := filepath.Base(pkg.ImportPath)
		output := name + ".apk"
		if len(buildO) > 0 {
			output = buildO
		}
		return Cmd(CmdContext{
			Verbose: true,
			ShowCmd: buildX,
		}, "adb", "install", "-r", output)
	case "ios":
		if err := runBuild(cmd); err != nil {
			return err
		}
		name := filepath.Base(pkg.ImportPath)
		output := name + ".app"
		if len(buildO) > 0 {
			output = buildO
		}
		return Cmd(CmdContext{
			Verbose: true,
			ShowCmd: buildX,
		}, "ios-deploy", "-b", output)
	}
	return fmt.Errorf("unsuported target: %s", buildTarget)
}
