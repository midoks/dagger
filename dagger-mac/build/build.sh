#!/bin/bash
# 打包

BUILD_DIR=$(cd "$(dirname "$0")"; pwd)
ROOT_DIR=$(dirname "$BUILD_DIR")

APP_NAME="dagger"
APP_VER=$(sed -n '/MARKETING_VERSION/{s/MARKETING_VERSION = //;s/;//;s/^[[:space:]]*//;p;q;}' $ROOT_DIR/dagger/dagger.xcodeproj/project.pbxproj)
DAGGER_RELEASE=${BUILD_DIR}/release

DMG_FINAL="${APP_NAME}.dmg"




function build(){
	mkdir -p $DAGGER_RELEASE

	echo "build dagger."${APP_VER}

	echo "Building archive... please wait a minute"
    xcodebuild -project $ROOT_DIR/dagger/dagger.xcodeproj -config Release -scheme dagger -archivePath ${DAGGER_RELEASE} archive

    echo "Exporting archive..."
    xcodebuild -archivePath ${BUILD_DIR}/release.xcarchive -exportArchive -exportPath ${DAGGER_RELEASE} -exportOptionsPlist $BUILD_DIR/build.plist

}


function createDmgByAppdmg(){

	# umount "/Volumes/${APP_NAME}"

	rm -rf ${BUILD_DIR}/${APP_NAME}.app ${BUILD_DIR}/${DMG_FINAL}
	\cp -Rf "${DAGGER_RELEASE}/${APP_NAME}.app" "${BUILD_DIR}/${APP_NAME}.app"

	# npm install -g appdmg
	echo ${BUILD_DIR}/appdmg.json
    appdmg appdmg.json ${DMG_FINAL}

    # umount "/Volumes/${APP_NAME}"
}

function makeDmg(){

	rm -fr ${DMG_FINAL} ${DAGGER_RELEASE}
	
	build
	createDmgByAppdmg
}


echo $ROOT_DIR
if [ "$1" = "build" ]
then
	build

else
    makeDmg
fi
echo 'done'
