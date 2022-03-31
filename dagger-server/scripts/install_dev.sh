#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin


TAGRT_DIR=/usr/local/dagger_dev
mkdir -p $TAGRT_DIR
cd $TAGRT_DIR


if [ ! -d $TAGRT_DIR/dagger ]; then
	git clone https://github.com/midoks/dagger
	cd $TAGRT_DIR/dagger
else
	cd $TAGRT_DIR/dagger
	git pull https://github.com/midoks/dagger
fi

cd $TAGRT_DIR/dagger/dagger-server

go mod tidy
go mod vendor


rm -rf dagger
go build ./


cd $TAGRT_DIR/dagger/dagger-server/scripts

sh make.sh

systemctl daemon-reload

service dagger restart

cd $TAGRT_DIR/dagger/dagger-server && ./dagger-server -v


