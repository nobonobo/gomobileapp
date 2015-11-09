package main

import (
	"encoding/xml"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

var (
	iconAttr = xml.Attr{
		Name: xml.Name{
			Space: "http://schemas.android.com/apk/res/android",
			Local: "icon",
		},
		Value: "@drawable/ic_launcher",
	}
)

func buildAndroid(pkg *build.Package, ctx CmdContext) error {
	name := filepath.Base(pkg.ImportPath)
	output := name + ".apk"
	if len(buildO) > 0 {
		output = buildO
	}
	subargs := append(ctx.Args, []string{
		"-target", buildTarget, "-o", output, pkg.ImportPath,
	}...)
	if err := Cmd(ctx, goMobileExe, subargs...); err != nil {
		return err
	}
	if len(customIcon) == 0 {
		return nil
	}

	_ = Cmd(ctx, "keytool", "-genkey",
		"-keystore", filepath.Join("build", "android", "app.keystore"),
		"-storepass", "secret",
		"-dname", fmt.Sprintf("CN=%s,O=apps.gomobile.org,C=JP", name),
		"-keypass", "secret",
		"-keyalg", "RSA", "-validity", "18250",
		"-alias", name,
	)

	_ = os.RemoveAll(filepath.Join("build", "android", "apk"))
	if err := Cmd(ctx, "apktool", "decode", "-o",
		filepath.Join("build", "android", "apk"), output); err != nil {
		return err
	}

	manifest := filepath.Join("build", "android", "apk", "AndroidManifest.xml")
	fp, err := os.Open(manifest)
	if err != nil {
		return err
	}
	defer fp.Close()
	var v *Tag
	if err := xml.NewDecoder(fp).Decode(&v); err != nil {
		return err
	}
	attrs := v.Attr
	v.Attr = []xml.Attr{}
	for _, a := range attrs {
		switch a.Name.Local {
		case "versionCode":
		case "versionName":
		default:
			v.Attr = append(v.Attr, a)
		}
	}

	var app *Tag
	for _, c := range v.Children {
		if cf, ok := c.(*Tag); ok {
			cf.Name.Local = "application"
			app = cf
			break
		}
	}
	if app == nil {
		return fmt.Errorf("<application> not found in AndroidManifest.xml")
	}
	app.Attr = append(app.Attr, iconAttr)
	fp.Close()
	fp, err = os.Create(manifest)
	if err != nil {
		return err
	}
	if err := xml.NewEncoder(fp).Encode(&v); err != nil {
		return err
	}

	_ = os.RemoveAll(filepath.Join("build", "android", "apk", "res"))
	if len(customIcon) > 0 {
		if err := Cmd(ctx, "icons", "--device", "android", "-o",
			filepath.Join("build", "android", "apk", "res"),
			customIcon,
		); err != nil {
			return err
		}
	}
	_ = os.Remove(filepath.Join("build", "android", "apk", "res", "playstore.png"))
	_ = os.RemoveAll(output)
	if err := Cmd(ctx, "apktool", "build", "-o", output,
		"-f", filepath.Join("build", "android", "apk")); err != nil {
		return err
	}
	if err := Cmd(ctx, "jarsigner", "-verbose",
		"-keystore", filepath.Join("build", "android", "app.keystore"),
		"-storepass", "secret",
		"-sigalg", "MD5withRSA",
		"-digestalg", "SHA1",
		"-tsa", "http://timestamp.digicert.com",
		output, name,
	); err != nil {
		return err
	}

	return nil
}
