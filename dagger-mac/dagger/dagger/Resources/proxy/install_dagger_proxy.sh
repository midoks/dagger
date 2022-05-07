#!/bin/sh

#  install_helper.sh
#  dagger

NEW_VERSION="0.0.2"

cd "$(dirname "${BASH_SOURCE[0]}")"


mkdir -p "$HOME/Library/Application\ Support/dagger"

F_VER=`$HOME/Library/Application\ Support/dagger/dagger-client-http -v | grep version | awk '{print $3}'`
echo $F_VER

if [ "$NEW_VERSION"!="$F_VER" ];then
    rm -rf "$HOME/Library/Application Support/dagger/dagger-client-http"
    cp dagger-client-http "$HOME/Library/Application Support/dagger/"
fi

launchctl stop com.midoks.dagger.http
launchctl unload "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"

