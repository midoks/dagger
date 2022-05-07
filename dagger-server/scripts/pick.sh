#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

# https://github.com/FiloSottile/homebrew-musl-cross
# brew install FiloSottile/musl-cross/musl-cross --without-x86_64 --with-i486 --with-aarch64 --with-arm

# brew install mingw-w64
# sudo port install mingw-w64

VERSION=0.0.5
curPath=`pwd`
rootPath=$(dirname "$curPath")

LDFLAGS="-w -s"
PACK_NAME=dagger-server

# go tool dist list
mkdir -p $rootPath/tmp/build
mkdir -p $rootPath/tmp/package

source ~/.bash_profile


echo $rootPath
cd $rootPath


echo $LDFLAGS
build_app(){

	if [ -f $rootPath/tmp/build/dagger-server ]; then
		rm -rf $rootPath/tmp/build/dagger-server
		rm -rf $rootPath/dagger-server
	fi

	if [ -f $rootPath/tmp/build/dagger-server.exe ]; then
		rm -rf $rootPath/tmp/build/dagger-server.exe
		rm -rf $rootPath/dagger-server.exe
	fi

	echo "build_app" $1 $2

	echo "export CGO_ENABLED=1 GOOS=$1 GOARCH=$2"
	echo "cd $rootPath && go build dagger.go"

	export CGO_ENABLED=1 GOOS=$1 GOARCH=$2
	# export CGO_ENABLED=1 GOOS=linux GOARCH=amd64


	if [ $1 != "darwin" ];then
		export CGO_ENABLED=1 GOOS=$1 GOARCH=$2
		export CGO_LDFLAGS="-static"
	fi

    #export CGO_LDFLAGS="-static"
	if [ $1 == "windows" ];then
		
		if [ $2 == "amd64" ]; then
			export CC=x86_64-w64-mingw32-gcc
			export CXX=x86_64-w64-mingw32-g++
		else
			export CC=i686-w64-mingw32-gcc
			export CXX=i686-w64-mingw32-g++
		fi

		cd $rootPath && go build -v -ldflags "${LDFLAGS}" -o dagger-server.exe
	fi

	if [ $1 == "linux" ]; then
		export CC=x86_64-linux-musl-gcc
		if [ $2 == "amd64" ]; then
			export CC=x86_64-linux-musl-gcc

		fi

		if [ $2 == "386" ]; then
			export CC=i486-linux-musl-gcc
		fi

		if [ $2 == "arm64" ]; then
			export CC=aarch64-linux-musl-gcc
		fi

		if [ $2 == "arm" ]; then
			export CC=arm-linux-musleabi-gcc
		fi

		cd $rootPath && go build -ldflags "${LDFLAGS}"
	fi

	if [ $1 == "darwin" ]; then
		echo "cd $rootPath && go build -v  -ldflags '${LDFLAGS}'"
		cd $rootPath && go build -v -ldflags "${LDFLAGS}"
	fi
	

	cp -r $rootPath/scripts $rootPath/tmp/build

	# cp -r $rootPath/conf $rootPath/tmp/build
	# sed 's/dev/prod/g' $rootPath/conf/app.conf > $rootPath/tmp/build/conf/app.conf

	mkdir -p $rootPath/tmp/build/logs
	echo "" > $rootPath/tmp/build/logs/README.md

	cd $rootPath/tmp/build && xattr -c * && rm -rf ./*/.DS_Store && rm -rf ./*/*/.DS_Store


	if [ $1 == "windows" ];then
		cp $rootPath/dagger-server.exe $rootPath/tmp/build
	else
		cp $rootPath/dagger-server $rootPath/tmp/build
	fi

	# tar.gz
	cd $rootPath/tmp/build && tar -zcvf ${PACK_NAME}_${VERSION}_$1_$2.tar.gz ./ && mv ${PACK_NAME}_${VERSION}_$1_$2.tar.gz $rootPath/tmp/package
	
}

golist=`go tool dist list`
echo $golist

build_app linux amd64
build_app linux 386
build_app linux arm64
build_app linux arm
build_app darwin amd64
build_app windows 386
build_app windows amd64

