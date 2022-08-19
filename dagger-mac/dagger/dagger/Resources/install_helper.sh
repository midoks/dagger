#!/bin/sh

#  install_helper.sh
#  dagger

cd "$(dirname "${BASH_SOURCE[0]}")"

sudo mkdir -p "/Library/Application Support/dagger/"

sudo cp dagger-helper "/Library/Application Support/dagger/"
sudo chown root:admin "/Library/Application Support/dagger/dagger-helper"
sudo chmod a+rx "/Library/Application Support/dagger/dagger-helper"
sudo chmod +s "/Library/Application Support/dagger/dagger-helper"

sudo cp dagger-cf "/Library/Application Support/dagger/"
sudo chown root:admin "/Library/Application Support/dagger/dagger-cf"
sudo chmod a+rx "/Library/Application Support/dagger/dagger-cf"
sudo chmod +s "/Library/Application Support/dagger/dagger-cf"

sudo cp ip.txt "/Library/Application Support/dagger"
sudo cp ipv6.txt "/Library/Application Support/dagger"
    
echo done
