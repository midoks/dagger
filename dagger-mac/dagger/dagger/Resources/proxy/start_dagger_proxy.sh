#!/bin/sh

#  install_helper.sh
#  dagger

NEW_VERSION="0.0.2"

F_VER=`$HOME/Library/Application\ Support/dagger/dagger-client-http -v | grep version | awk '{print $3}'`
echo $F_VER

if [ "$NEW_VERSION"!="$F_VER" ];then
    rm -rf "$HOME/Library/Application\ Support/dagger/dagger-client-http"
    cp dagger-client-http "$HOME/Library/Application Support/dagger/"
fi

echo "" > "$HOME/Library/Logs/dagger-client-http.log"

chmod 644 "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl load -wF "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl start com.midoks.dagger.http

