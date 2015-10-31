# gomobileapp

cli wrapper with additional features for gomobile.

## features

- icons generation from one image.
- add icon resource of app. (https://github.com/golang/go/issues/9985)
- run on desktop
- install app for iOS

## install

```sh
go get -u golang.org/x/mobile/cmd/...
go get github.com/nobonobo/gomobileapp
gomobileapp init
```

## prerequire

### for iOS

- Xcode7 install
- npm install -g ios-deploy

**Create Certificate for iOS Code signing**

- Open Preferences of Xcode7.
- Select Accounts and append own account.
- Click 'View Details' button for own account.
- Show 'Sigining Identities'
- Click 'Create' button for 'iOS Development'

### for Android

- android-platform-tools install
- android-apktool install
- java runtime environment(> 1.8) install

### icon generator

```sh
curl -kL https://raw.github.com/pypa/pip/master/contrib/get-pip.py | sudo python2
sudo pip2 install icons
```

## usage

for iOS app
```sh
go get -d github.com/nobonobo/nobopiano
gomobileapp build -icon icon.png -target ios github.com/nobopiano/nobopiano
ios-deploy install -r nobopiano.apk
```

for Android apk
```sh
go get -d github.com/nobonobo/nobopiano
gomobileapp build -icon icon.png -target android github.com/nobopiano/nobopiano
adb install -r nobopiano.apk
```

for Desktop app
```sh
go get -d github.com/nobonobo/nobopiano
gomobileapp build -target desktop github.com/nobopiano/nobopiano
./nobopiano
```

## TODO

- no dependency icons.
- support install -target ios
- support icon for desktop
