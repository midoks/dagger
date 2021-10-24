#!/bin/sh

#  install_helper.sh
#  dagger

cd "$(dirname "${BASH_SOURCE[0]}")"

sudo mkdir -p "/Library/Application Support/dagger/"
sudo cp dagger-helper "/Library/Application Support/dagger/"
sudo chown root:admin "/Library/Application Support/dagger/dagger-helper"
sudo chmod a+rx "/Library/Application Support/dagger/dagger-helper"
sudo chmod +s "/Library/Application Support/dagger/dagger-helper"

echo done
