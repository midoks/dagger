#!/bin/sh

#  install_helper.sh
#  dagger

echo "" > "$HOME/Library/Logs/dagger-client-http.log"

chmod 644 "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl load -wF "$HOME/Library/LaunchAgents/com.midoks.dagger.http.plist"
launchctl start com.midoks.dagger.http

