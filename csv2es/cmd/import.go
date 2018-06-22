// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ibbd-dev/go-csv"
	"github.com/spf13/cobra"
	"gopkg.in/olivere/elastic.v5"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import data from csv to es",
	Long: `import data from csv to es
`,
	Example: `
csv2es import --host=locahost --port=9200 --mapping=mapping_filename.json --index=test --type=test --csv=source.csv
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 创建输入文件
		in, err := os.Open(cParams.CsvFilename)
		if err != nil {
			panic(err)
		}
		defer in.Close()

		reader := goCsv.NewMapReader(in)
		fieldnames, err := reader.GetFieldnames()
		if err != nil {
			panic(err)
		}
		fmt.Printf("fieldnames: %s\n", strings.Join(fieldnames, ", "))

		// es
		conn, err := getESConnect()
		ctx := context.Background()
		exists, err := conn.IndexExists(cParams.IndexName).Do(ctx)
		if err != nil {
			panic(fmt.Errorf("check index exists error: %v", err.Error()))
		}

		if exists {
			if cParams.DeleteIndex {
				// 删除同名索引
				fmt.Printf("begin to delete index: %s\n", cParams.IndexName)
				if _, err = conn.DeleteIndex(cParams.IndexName).Do(ctx); err != nil {
					panic(fmt.Errorf("delete index error: %s", cParams.IndexName))
				}
			}

			// 创建索引
			fmt.Printf("begin to create index: %s\n", cParams.IndexName)
			if _, err = conn.CreateIndex(cParams.IndexName).Do(ctx); err != nil {
				panic(err)
			}

			if len(cParams.Mapping) > 0 {
				bytes, err := ioutil.ReadFile(cParams.Mapping)
				if err != nil {
					fmt.Println("Read mapping file: ", err.Error())
					panic(err)
				}

				var mapping = make(map[string]interface{})
				if err := json.Unmarshal(bytes, &mapping); err != nil {
					fmt.Println("Json Unmarshal: ", err.Error())
					panic(err)
				}

				fmt.Printf("begin to put index mapping: %s\n", cParams.IndexName)
				if _, err := conn.PutMapping().Index(cParams.IndexName).Type(cParams.DocType).BodyJson(mapping).Do(ctx); err != nil {
					panic(err)
				}
			}
		}

		var count int
		bulk := conn.Bulk()
		for {
			row, err := reader.Read()
			if err == io.EOF {
				fmt.Println("the file is over")
				break
			}
			if err != nil {
				panic(err)
			}

			req := elastic.NewBulkIndexRequest().Index(cParams.IndexName).Type(cParams.DocType).Doc(row)
			bulk.Add(req)

			count++
			if cParams.Size > 0 && count >= cParams.Size {
				fmt.Printf("导出限制数量：%d\n", cParams.Size)
				break
			}
		}

		// 执行导入
		bulkResponse, err := bulk.Do(ctx)
		if err != nil {
			panic(fmt.Errorf("index %v 批量导入数据出错: %v", cParams.IndexName, err.Error()))
		}

		// 统计写入状态
		var errCount int
		indexed := bulkResponse.Indexed()
		for i, res := range indexed {
			if res.Error != nil {
				fmt.Printf("ERROR: %d, %+v\n", i, res.Error)
				errCount++
			}
		}
		fmt.Printf("向es写入的数据量：%d，异常：%d\n", count, errCount)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")
	importCmd.PersistentFlags().BoolVar(&cParams.DeleteIndex, "delete-index", false, "delete the same index before")
	importCmd.PersistentFlags().StringVar(&cParams.Mapping, "mapping", "", "the mapping file name, json format")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
