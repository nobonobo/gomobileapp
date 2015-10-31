package main

import (
	"go/build"
	"path/filepath"
)

func buildDesktop(pkg *build.Package, ctx CmdContext) error {
	output := filepath.Base(pkg.ImportPath)
	if len(buildO) > 0 {
		output = buildO
	}
	subargs := append(ctx.Args, []string{
		"-o", output, pkg.ImportPath,
	}...)
	if err := Cmd(ctx, "go", subargs...); err != nil {
		return err
	}
	_ = Cmd(ctx, "cp", "-rf",
		filepath.Join(pkg.Dir, "assets"),
		filepath.Dir(output),
	)
	return nil
}
