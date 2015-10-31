package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildEnviOS(arch string) (map[string]string, error) {
	env := map[string]string{}
	m := func(k string) string { return env[k] }
	env["GOMOBILE"] = filepath.Clean(filepath.Join(filepath.Dir(goMobileExe), "..", "pkg", "gomobile"))
	env["GOOS"] = "darwin"
	env["GOARCH"] = arch
	env["CC"] = "/Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/clang"
	env["CXX"] = env["CC"]
	env["SDK"] = "/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/Developer/SDKs/iPhoneOS.sdk"
	env["CGO_ENABLED"] = "1"
	switch arch {
	case "arm":
		env["GOARM"] = "7"
		env["CGO_CFLAGS"] = os.Expand("-isysroot ${SDK} -arch armv7", m)
		env["CGO_LDFLAGS"] = os.Expand("-isysroot ${SDK} -arch armv7", m)
		env["PKGDIR"] = os.Expand("${GOMOBILE}/pkg_darwin_arm", m)
	case "arm64":
		env["CGO_CFLAGS"] = os.Expand("-isysroot ${SDK} -arch arm64", m)
		env["CGO_LDFLAGS"] = os.Expand("-isysroot ${SDK} -arch arm64", m)
		env["PKGDIR"] = os.Expand("${GOMOBILE}/pkg_darwin_arm64", m)
	default:
		return nil, fmt.Errorf("unknown arch: %q", arch)
	}
	return env, nil
}

func buildiOS(pkg *build.Package, ctx CmdContext) error {
	var firstErr error
	output := filepath.Base(pkg.ImportPath) + ".app"
	if len(buildO) > 0 {
		output = buildO
	}
	if _, err := os.Stat(filepath.Join("build", "ios")); err != nil {
		subargs := []string{
			"build", "-target", "ios", "-o", output, "-work", pkg.ImportPath,
		}
		c := exec.Command(goMobileExe, subargs...)
		buff := bytes.NewBuffer(nil)
		c.Stdout = buff
		c.Stderr = os.Stderr
		firstErr = c.Run()
		if firstErr != nil {
			log.Println(firstErr)
		}
		b, _, _ := bufio.NewReader(buff).ReadLine()
		def := string(b)
		if !strings.HasPrefix(def, "WORK=") {
			return fmt.Errorf("not found working folder")
		}
		work := strings.TrimPrefix(def, "WORK=")
		if err := os.MkdirAll("build", 0755); err != nil {
			return err
		}
		os.RemoveAll(filepath.Join("build", "ios"))
		if err := os.Rename(work, filepath.Join("build", "ios")); err != nil {
			return err
		}
	}
	if firstErr != nil {
		Cmd(ctx, "open", filepath.Join("build", "ios", "main.xcodeproj"))
	}
	if len(customIcon) > 0 {
		if err := Cmd(ctx, "icons", "--device", "ios", "-o",
			filepath.Join("build", "ios", "main", "Images.xcassets"),
			customIcon,
		); err != nil {
			return err
		}
	}
	buildBin := func(arch string) error {
		envs, err := buildEnviOS(arch)
		if err != nil {
			return err
		}
		envs["WORK"] = filepath.Join("build", "ios")
		m := func(k string) string { return envs[k] }
		ctx.Env = map2envs(envs)
		subargs := []string{
			"build", "-p=4",
			os.Expand("-pkgdir=${PKGDIR}", m),
			"-tags=ios",
			os.Expand("-o=${WORK}/${GOARCH}", m),
		}
		subargs = append(subargs, ctx.Args[1:]...)
		err = Cmd(ctx, "go", append(subargs, pkg.ImportPath)...)
		if err != nil {
			return err
		}
		return nil
	}
	if err := buildBin("arm"); err != nil {
		return err
	}
	if err := buildBin("arm64"); err != nil {
		return err
	}
	envs := map[string]string{
		"WORK": filepath.Join("build", "ios"),
	}
	m := func(k string) string { return envs[k] }
	if err := Cmd(ctx, "xcrun", "lipo", "-create",
		os.Expand("${WORK}/arm", m),
		os.Expand("${WORK}/arm64", m),
		"-o",
		os.Expand("${WORK}/main/main", m),
	); err != nil {
		return err
	}
	envs["PKGDIR"] = pkg.Dir
	_ = os.RemoveAll(os.Expand("${WORK}/main/assets", m))
	if err := Cmd(ctx,
		"cp", "-Rf",
		os.Expand("${PKGDIR}/assets", m),
		os.Expand("${WORK}/main/", m),
	); err != nil {
		return err
	}
	if err := Cmd(ctx,
		"xcrun", "xcodebuild", "-configuration", "Release", "-project",
		os.Expand("${WORK}/main.xcodeproj", m),
	); err != nil {
		return err
	}
	_ = os.RemoveAll(output)
	//mv $WORK/build/Release-iphoneos/main.app $OUTPUT
	if err := os.Rename(
		os.Expand("${WORK}/build/Release-iphoneos/main.app", m),
		output,
	); err != nil {
		return err
	}
	return nil
}
