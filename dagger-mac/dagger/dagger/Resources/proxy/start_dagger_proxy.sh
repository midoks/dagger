#!/bin/sh

#  install_helper.sh
#  dagger

chmod 644 "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl load -wF "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl start com.midoks.dagger.http

