#!/bin/bash
# 打包

BUILD_DIR=$(cd "$(dirname "$0")"; pwd)
ROOT_DIR=$(dirname "$BUILD_DIR")

APP_NAME="dagger"
APP_VER=$(sed -n '/MARKETING_VERSION/{s/MARKETING_VERSION = //;s/;//;s/^[[:space:]]*//;p;q;}' $ROOT_DIR/dagger/dagger.xcodeproj/project.pbxproj)
DAGGER_RELEASE=${BUILD_DIR}/release

mkdir -p $DAGGER_RELEASE


function build(){
	echo "build dagger."${APP_VER}

	echo "Building archive... please wait a minute"
    xcodebuild -project $ROOT_DIR/dagger/dagger.xcodeproj -config Release -scheme dagger -archivePath ${DAGGER_RELEASE} archive

    echo "Exporting archive..."
    xcodebuild -archivePath ${DAGGER_RELEASE} -exportArchive -exportPath ${DAGGER_RELEASE} -exportOptionsPlist ./build.plist


}


echo $ROOT_DIR

build