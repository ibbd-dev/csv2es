#!/bin/bash
# 
# csv数据导入测试脚本
# Author: alex
# Created Time: 2018年06月21日 星期四 18时46分23秒

cmd="csv2es"
if [ -f csv2es ]; then
    cmd="./csv2es"
fi

$cmd import --host=100.115.147.50 --port=9200 --index=test --type=test --mapping=./eyenlp_area2016.json --csv=./eyenlp_area2016.csv --delete-index=true --debug=true

# 查看mapping是否生效
curl 100.115.147.50:9200/test/_mapping|json_pp
