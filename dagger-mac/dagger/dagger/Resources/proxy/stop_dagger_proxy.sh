#!/bin/sh

#  install_helper.sh
#  dagger

launchctl stop com.midoks.dagger.http.plist
launchctl unload "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"

