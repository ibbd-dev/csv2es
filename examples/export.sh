#!/bin/bash
# 
# 数据导出
# Author: alex
# Created Time: 2018年06月21日 星期四 21时21分19秒

cmd="csv2es"
if [ -f csv2es ]; then
    cmd="./csv2es"
fi

index=test
if [ $# = 1 ]; then
    index="$1"
fi

$cmd export --host=100.115.147.50 --port=9200 --index="$index" --csv=./output.csv --bulk-size=300 --debug=true

