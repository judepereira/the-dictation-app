#!/bin/bash

set -euo pipefail

ROOT_DIR="build/TheDictationApp.app"

CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ${ROOT_DIR}/Contents/MacOS/TheDictationApp dictation/cmd/app

#CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin-amd64 ./cmd/app
#lipo -create -output Dictation.app/Contents/MacOS/Dictation Dictation.app/Contents/MacOS/Dictation bin-amd64

# Bundle structure
mkdir -p ${ROOT_DIR}/Contents/{MacOS,Resources}
cp assets/icon.png ${ROOT_DIR}/Contents/Resources/AppIcon.png
#cp -R models/ ${ROOT_DIR}/Contents/Resources/models  # include e.g., base.en.bin
plutil -convert xml1 -o ${ROOT_DIR}/Contents/Info.plist Info.plist

#codesign --deep --force --options runtime \
#  --entitlements entitlements.plist \
#  -s "Developer ID Application: Your Name (TEAMID)" \
#  ${ROOT_DIR}

#xcrun notarytool submit Dictation.app --keychain-profile "notaryprofile" --wait
#xcrun stapler staple Dictation.app


hdiutil create -volname "The Dictation App" -srcfolder ${ROOT_DIR} -ov -format UDZO build/TheDictationApp.dmg

echo "Build complete!"