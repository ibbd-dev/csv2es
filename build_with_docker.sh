#!/bin/bash
# 
# 在本地环境编译所有工具
# Author: alex
# Created Time: 2017年07月05日 星期三 11时47分58秒


# 判断容器是否已经创建
if 
    sudo docker ps -a|grep 'ibbd-xgo'; 
then
    echo "ibbd-xgo is exist"
else
    # 容器还没创建
    echo 'begin to run ibbd-xgo:'
    echo "可能需要手动启动"
    if 
        sudo docker run -ti --name=ibbd-xgo -v /var/www/:/var/www/ karalabe/xgo-latest /bin/bash
    then
        echo 'build success!'
        exit 0
    else
        echo 'build failure!'
        exit 1
    fi
fi

# 判断容器是否正在运行
if 
    sudo docker ps|grep 'ibbd-xgo'
then
    echo 'ibbd-xgo is running!'
else
    if 
        sudo docker start ibbd-xgo; 
    then
        echo 'ibbd-xgo start success!'
    else
        echo 'ibbd-xgo start failure!'
        exit 1
    fi
fi

if [ ! -d /var/www/build/elasticsearch ]; then
    mkdir -p /var/www/build/elasticsearch
fi


version=`cat ./csv2es/cmd/root.go|grep "Version: "|cut -d'"' -f2`
if [ ${#version} = 0 ]; then 
    echo "version is empty"
    exit 1
fi

echo "current version: $version"

# 执行编译
curr_path=$PWD
echo "Begin compile:"
sudo docker exec -ti ibbd-xgo "$curr_path/build_in_xgo.sh" $version

user=`whoami`
sudo chown -R $user:$user /var/www/build/elasticsearch

# 处理相关文件
path=/var/www/build/elasticsearch/$version
if [ ! -d $path ]; then
    echo "$path is not existed!"
    exit 1
fi
path=/var/www/build/elasticsearch/$version/examples
if [ ! -d $path ]; then
    mkdir $path
fi
cp ./examples/* $path/

# 本地编译
go get github.com/ibbd-dev/csv2es/csv2es

# 
cd csv2es
go build
rm -f csv2es-$version
mv csv2es csv2es-$version
