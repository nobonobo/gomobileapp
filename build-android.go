package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	manifest = template.Must(template.New("").Parse(
		`<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android" package="org.golang.todo.{{.NAME}}" android:versionCode="1" android:versionName="1.0">
  <uses-sdk android:minSdkVersion="9" />
  <application android:label="{{.LABEL}}" android:debuggable="true" android:icon="@drawable/ic_launcher">
    <activity android:name="org.golang.app.GoNativeActivity" android:label="{{.LABEL}}" android:configChanges="orientation|keyboardHidden">
      <meta-data android:name="android.app.lib_name" android:value="{{.NAME}}" />
      <intent-filter>
        <action android:name="android.intent.action.MAIN" />
        <category android:name="android.intent.category.LAUNCHER" />
      </intent-filter>
    </activity>
  </application>
</manifest>
`,
	))
)

func buildAndroid(pkg *build.Package, ctx CmdContext) error {
	name := filepath.Base(pkg.ImportPath)
	label := strings.ToUpper(name[:1]) + name[1:]
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

	mf, err := os.Create(filepath.Join("build", "android", "apk", "AndroidManifest.xml"))
	if err != nil {
		return err
	}
	if err := manifest.Execute(mf, map[string]interface{}{
		"NAME":  name,
		"LABEL": label,
	}); err != nil {
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
