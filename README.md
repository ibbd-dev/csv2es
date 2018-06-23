# csv2es: the data import/export tool between es and csv
elasticsearch与csv之间的数据导入导出工具

说明：

- 现在版本只支持`elasticsearch 5.*`版本
- 该工具只是实现对数据的导入导出，并不包含相关的ETL过程。

## Install 

```sh
go get -u github.com/ibbd-dev/csv2es/csv2es
```

## Example

```sh
# help 
csv2es --help

# import data to es from csv
csv2es import --host=locahost --port=9200 --mapping=mapping_filename.json --index=test --type=test --csv=source.csv

# export data to csv from es
csv2es export --host=locahost --port=9200 --index=test --csv=source.csv
```

