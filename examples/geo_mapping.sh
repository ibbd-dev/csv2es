#!/bin/bash
# 
# es geo_point mapping test
# Author: alex
# Created Time: 2018年06月22日 星期五 20时06分28秒
curl -XDELETE 100.115.147.50:9200/my_index

curl -XPUT 100.115.147.50:9200/my_index -d '{
  "mappings": {
    "my_doc": {
      "properties": {
        "location": {
          "type": "geo_point"
        }
      }
    }
  }
}'

curl -XPUT 100.115.147.50:9200/my_index/my_doc/1 -d '{
  "text": "Geo-point as a string",
  "location": "41.12,-71.34"
}'

curl -XPUT 100.115.147.50:9200/my_index/my_doc/2 -d '{
  "text": "Geo-point as a map",
  "location": {
    "lat": 41.12,
    "lon": -71.34
  }
}'

curl 100.115.147.50:9200/my_index/my_doc/1 |json_pp
curl 100.115.147.50:9200/my_index/my_doc/2 |json_pp
