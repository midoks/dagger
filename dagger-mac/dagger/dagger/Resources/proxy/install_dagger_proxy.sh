#!/bin/sh

#  install_helper.sh
#  dagger

cd "$(dirname "${BASH_SOURCE[0]}")"

mkdir -p "$HOME/Library/Application Support/dagger"
cp dagger-client-http "$HOME/Library/Application Support/dagger/"

launchctl stop com.midoks.dagger.http
launchctl unload "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"

