// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// ctx, pkg, tmpdir in build.go

var cmdBind = &command{
	run:   runBind,
	Name:  "bind",
	Usage: "[-target android|ios] [-o output] [build flags] [package]",
	Short: "build a library for Android and iOS",
	Long: `
Bind generates language bindings for the package named by the import
path, and compiles a library for the named target system.

The -target flag takes a target system name, either android (the
default) or ios.

For -target android, the bind command produces an AAR (Android ARchive)
file that archives the precompiled Java API stub classes, the compiled
shared libraries, and all asset files in the /assets subdirectory under
the package directory. The output is named '<package_name>.aar' by
default. This AAR file is commonly used for binary distribution of an
Android library project and most Android IDEs support AAR import. For
example, in Android Studio (1.2+), an AAR file can be imported using
the module import wizard (File > New > New Module > Import .JAR or
.AAR package), and setting it as a new dependency
(File > Project Structure > Dependencies).  This requires 'javac'
(version 1.7+) and Android SDK (API level 9 or newer) to build the
library for Android. The environment variable ANDROID_HOME must be set
to the path to Android SDK. The generated Java class is in the java
package 'go.<package_name>' unless -javapkg flag is specified.

For -target ios, gomobile must be run on an OS X machine with Xcode
installed. Support is not complete. The generated Objective-C types
are prefixed with 'Go' unless the -prefix flag is provided.

The -v flag provides verbose output, including the list of packages built.

The build flags -a, -i, -n, -x, -gcflags, -ldflags, -tags, and -work
are shared with the build command. For documentation, see 'go help build'.
`,
}

func runBind(cmd *command) error {
	return Cmd(CmdContext{}, goMobileExe, append([]string{"bind"}, cmd.flag.Args()...)...)
}
