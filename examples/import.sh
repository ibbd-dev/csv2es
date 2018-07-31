#!/bin/bash
# 
# csv数据导入测试脚本
# Author: alex
# Created Time: 2018年06月21日 星期四 18时46分23秒

cmd="csv2es"
if [ -f csv2es ]; then
    cmd="./csv2es"
fi

host=$1

$cmd import --host=$host --port=9200 --index=test_1 --type=test_1 --mapping=./eyenlp_area2016.json --csv=./eyenlp_area2016.csv --delete-index=true --debug=true

$cmd import --host=$host --port=9200 --index=test_2 --type=test_2 --mapping=./eyenlp_housing_estate_20180523.json --csv=./eyenlp_housing_estate_20180523.csv --delete-index=true --debug=true

# 查看mapping是否生效
curl $host:9200/test_1/_mapping|json_pp
curl $host:9200/test_2/_mapping|json_pp

