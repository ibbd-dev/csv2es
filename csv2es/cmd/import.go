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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ibbd-dev/go-csv"
	"github.com/ibbd-dev/go-tools/es"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import data from csv to es",
	Long: `从csv导入数据到es
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
		conn, err := es.NewClient(cParams.Host, cParams.Port, cParams.IndexName, cParams.DocType)
		if err != nil {
			panic(fmt.Errorf("es newClient: %v", err.Error()))
		}

		if cParams.BulkSize <= 0 {
			cParams.BulkSize = 1000
		}
		conn.SetLimit(cParams.Limit)
		conn.SetBulkSize(cParams.BulkSize)
		conn.SetDebug(cParams.Debug)

		var mapping = make(map[string]interface{})
		if len(cParams.Mapping) > 0 {
			bytes, err := ioutil.ReadFile(cParams.Mapping)
			if err != nil {
				fmt.Println("Read mapping file: ", err.Error())
				panic(err)
			}

			if err := json.Unmarshal(bytes, &mapping); err != nil {
				fmt.Println("Json Unmarshal: ", err.Error())
				panic(err)
			}
		}

		if err = conn.ImportInit(cParams.DeleteIndex, mapping); err != nil {
			panic(err)
		}

		var count int
		for {
			row, err := reader.Read()
			if err == io.EOF {
				fmt.Println("the file is over")
				break
			}
			if err != nil {
				panic(err)
			}

			conn.BulkAdd(row)

			count++
			if cParams.Limit > 0 && count >= cParams.Limit {
				fmt.Printf("导出限制数量：%d\n", cParams.Limit)
				break
			}
		}

		// 执行导入
		if err = conn.BulkImport(); err != nil {
			panic(err)
		}
		fmt.Printf("向es写入的数据量：%d\n", count)
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().BoolVar(&cParams.DeleteIndex, "delete-index", false, "delete the same index before")
	importCmd.PersistentFlags().StringVar(&cParams.Mapping, "mapping", "", "the mapping file name, json format")
}
