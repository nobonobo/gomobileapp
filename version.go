// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var cmdVersion = &command{
	run:   runVersion,
	Name:  "version",
	Usage: "",
	Short: "print version",
	Long: `
Version prints versions of the gomobile binary and tools
`,
}

func runVersion(cmd *command) (err error) {
	return Cmd(CmdContext{}, goMobileExe, "version")
}
