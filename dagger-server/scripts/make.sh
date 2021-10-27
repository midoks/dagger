#!/bin/bash

_os=`uname`
_path=`pwd`
_dir=`dirname $_path`

sed "s:{APP_PATH}:${_dir}:g" $_dir/scripts/init.d/dagger.tpl > $_dir/scripts/init.d/dagger
chmod +x $_dir/scripts/init.d/dagger


if [ -d /etc/init.d ];then
	cp $_dir/scripts/init.d/dagger /etc/init.d/dagger
	chmod +x /etc/init.d/dagger
fi

echo `dirname $_path`