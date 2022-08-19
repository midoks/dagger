#!/bin/sh

#  install-dagger-cf.sh
#  dagger
#

cd "$(dirname "${BASH_SOURCE[0]}")"

NEW_VERSION="0.0.1"

F_VER=`$HOME/Library/Application\ Support/dagger/dagger-cf -v | grep version | awk '{print $3}'`
echo $F_VER

mkdir -p "$HOME/Library/Application\ Support/dagger"
if [ "$NEW_VERSION"!="$F_VER" ];then

    rm -rf "$HOME/Library/Application Support/dagger/dagger-cf"
    cp -rf dagger-cf "$HOME/Library/Application Support/dagger/"
    cp -rf ip.txt "$HOME/Library/Application\ Support/dagger"
    cp -rf ipv6.txt "$HOME/Library/Application\ Support/dagger"
    
    #sudo chown root:admin "/Library/Application Support/dagger/dagger-cf"
    #sudo chmod a+rx "/Library/Application Support/dagger/dagger-cf"
    #sudo chmod +s "/Library/Application Support/dagger/dagger-cf"
fi
echo done
