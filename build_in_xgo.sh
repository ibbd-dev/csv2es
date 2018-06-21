#!/bin/bash
# 
# 在xgo容器内部编译所有
# Author: alex
# Created Time: 2017年07月05日 星期三 11时43分21秒

export GOPATH=/var/www/

# 获取版本号
version=$1
targets="linux/amd64,windows/*,darwin/*"
package=github.com/ibbd-dev/csv2es/csv2es

xgo -out csv2es-"$version" --targets="$targets" "$package"

path=/var/www/build/elasticsearch/$version
if [ ! -d $path ]; then
    mkdir -p $path
fi
mv /build/* $path

echo 
echo "Finish!"
