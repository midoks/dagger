#!/bin/sh

#  install_helper.sh
#  dagger

cd "$(dirname "${BASH_SOURCE[0]}")"

echo $(dirname "${BASH_SOURCE[0]}")
mkdir -p "$HOME/Library/Application Support/dagger/"
cp dagger-client-http "$HOME/Library/Application Support/dagger/"

echo done
