package main

import (
	"fmt"
	"go/build"
	"os"
)

var cmdBuild = &command{
	run:   runBuild,
	Name:  "build",
	Usage: "[-target android|ios] [-icon icon.png] [-o output] [build flags] [package]",
	Short: "compile android APK and iOS app",
	Long: `
Build compiles and encodes the app named by the import path.

The named package must define a main function.

The -target flag takes a target system name, either android (the
default) or ios.

For -target android, if an AndroidManifest.xml is defined in the
package directory, it is added to the APK output. Otherwise, a default
manifest is generated.

For -target ios, gomobile must be run on an OS X machine with Xcode
installed. Support is not complete.

If the package directory contains an assets subdirectory, its contents
are copied into the output.

The -o flag specifies the output file name. If not specified, the
output file name depends on the package built.

The -v flag provides verbose output, including the list of packages built.

The build flags -a, -i, -n, -x, -gcflags, -ldflags, -tags, and -work are
shared with the build command. For documentation, see 'go help build'.

The -icon flag specifies the png image filename.
`,
}

func runBuild(cmd *command) (err error) {
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

	sub := []string{"build"}
	if len(buildGcflags) > 0 {
		sub = append(sub, "-gcflags", buildGcflags)
	}
	if len(buildLdflags) > 0 {
		sub = append(sub, "-ldflags", buildLdflags)
	}
	if buildA {
		sub = append(sub, "-a")
	}
	if buildI {
		sub = append(sub, "-i")
	}
	if buildN {
		sub = append(sub, "-n")
	}
	if buildV {
		sub = append(sub, "-v")
	}
	if buildX {
		sub = append(sub, "-x")
	}
	if buildWork {
		sub = append(sub, "-work")
	}
	switch buildTarget {
	case "android":
		if err := buildAndroid(pkg, CmdContext{
			Verbose: buildV,
			ShowCmd: buildX,
			Args:    sub,
		}); err != nil {
			return err
		}
	case "ios":
		if err := buildiOS(pkg, CmdContext{
			Verbose: buildV,
			ShowCmd: buildX,
			Args:    sub,
		}); err != nil {
			return err
		}
	case "desktop":
		if err := buildDesktop(pkg, CmdContext{
			Verbose: buildV,
			ShowCmd: buildX,
			Args:    sub,
		}); err != nil {
			return err
		}
	default:
		return fmt.Errorf(`unknown -target, %q.`, buildTarget)
	}
	return nil
}

func map2envs(m map[string]string) []string {
	res := []string{}
	for k, v := range m {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res
}
