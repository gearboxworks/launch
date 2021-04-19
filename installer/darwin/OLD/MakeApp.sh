#!/bin/bash

# Mac OSX .app builder

function die {
	echo "ERROR: $1" > /dev/null 1>&2
	exit 1
}

#if [ "$#" -ne 2 ]; then
#	die "Usage: `basename $0` AppNameHere icon-file.svg"
#fi

APPNAME=Launch
ICONNAME=./Gearbox.png

SAVED_DIR="${PWD}"

mkdir -p "app/$APPNAME.app/Contents/"{MacOS,Resources}

cat > "app/$APPNAME.app/Contents/Info.plist" <<END
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleGetInfoString</key>
  <string>$APPNAME</string>
  <key>CFBundleExecutable</key>
  <string>$APPNAME</string>
  <key>CFBundleIdentifier</key>
  <string>com.gearboxworks.launch</string>
  <key>CFBundleName</key>
  <string>$APPNAME</string>
  <key>CFBundleIconFile</key>
  <string>icon.icns</string>
  <key>CFBundleShortVersionString</key>
  <string>0.01</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>IFMajorVersion</key>
  <integer>0</integer>
  <key>IFMinorVersion</key>
  <integer>1</integer>
  <key>NSHighResolutionCapable</key><true/>
  <key>NSSupportsAutomaticGraphicsSwitching</key><true/>
</dict>
</plist>
END

cp $ICONNAME "app/$APPNAME.app/Contents/Resources/"
cd "app/$APPNAME.app/Contents/Resources/"

fileName="$(basename $ICONNAME)"
postfix=${fileName##*.}

if [[ $postfix == 'svg' ]]; then
    qlmanage -z -t -s 1024 -o ./ "$fileName"
    fileName=${fileName}.png
fi

echo $fileName

mkdir icon.iconset

sips -z 16 16 "$fileName" --out icon.iconset/icon_16x16.png
sips -z 32 32 "$fileName" --out icon.iconset/icon_16x16@2x.png
cp icon.iconset/icon_16x16@2x.png icon.iconset/icon_32x32.png
sips -z 64 64 "$fileName" --out icon.iconset/icon_32x32@2x.png
sips -z 128 128 "$fileName" --out icon.iconset/icon_128x128.png
sips -z 256 256 "$fileName" --out icon.iconset/icon_128x128@2x.png
cp icon.iconset/icon_128x128@2x.png icon.iconset/icon_256x256.png
sips -z 512 512 "$fileName" --out icon.iconset/icon_256x256@2x.png
cp icon.iconset/icon_256x256@2x.png icon.iconset/icon_512x512.png
sips -z 1024 1024 "$fileName" --out icon.iconset/icon_512x512@2x.png

# Create .icns file
iconutil -c icns icon.iconset

# Cleanup
rm -R icon.iconset
rm $fileName

cd "${SAVED_DIR}"
go run MakeDMG.go \
	-assets dmg/ \
	-bin ./launch \
	-icon ${ICONNAME} \
	-identifier com.gearboxworks.launch \
	-name "${APPNAME}" \
	-dmg "${APPNAME}.dmg" \
	-o ./dmg


