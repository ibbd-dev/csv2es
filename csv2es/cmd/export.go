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
	"os"
	"strings"

	"github.com/ibbd-dev/go-csv"
	"github.com/ibbd-dev/go-tools/es"
	"github.com/spf13/cobra"
	"gopkg.in/olivere/elastic.v5"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export data to csv from es",
	Long: `从es导出数据到csv文件,
`,
	Example: `
csv2es export --host=locahost --port=9200 --index=test --csv=source.csv
`,
	Run: func(cmd *cobra.Command, args []string) {
		// 创建输出文件
		out, err := os.Create(cParams.CsvFilename)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		writer := goCsv.NewMapWriterSimple(out)

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
		if cParams.Debug {
			fmt.Printf("client config: %+v\n", conn)
		}

		var query elastic.Query
		if len(cParams.QueryField) > 0 {
			query = elastic.NewTermQuery(cParams.QueryField, cParams.QueryValue)
		}
		if err = conn.SearchInit(query); err != nil {
			panic(err)
		}

		var count int // 记录总的记录数
		for {
			row, err := conn.Read()
			if err == io.EOF {
				fmt.Println("read over")
				break
			}
			if err != nil {
				panic(err)
			}

			if count == 0 {
				// 首行
				var headers []string
				for k, _ := range row {
					headers = append(headers, k)
				}

				writer.SetHeader(headers)
				if err = writer.WriteHeader(); err != nil {
					panic(fmt.Errorf("csv writer header error: %s", err.Error()))
				}
				fmt.Printf("Fieldnames: %s\n", strings.Join(headers, ", "))
			}

			var strRow = make(map[string]string)
			for k, v := range row {
				if vv, ok := v.(string); ok {
					strRow[k] = vv
				} else if sv, err := json.Marshal(v); err != nil {
					panic(fmt.Errorf("csv writer header error: %s", err.Error()))
				} else {
					strRow[k] = string(sv)
				}
			}

			count += 1
			writer.WriteRow(strRow)
			if count%cParams.BulkSize == 0 {
				writer.Flush()
			}
		} // end of for

		writer.Flush()
		fmt.Printf("Total %d\n", count)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.PersistentFlags().StringVar(&cParams.QueryField, "query-field", "", "过滤字段")
	exportCmd.PersistentFlags().StringVar(&cParams.QueryValue, "query-value", "", "过滤字段对应的值")
}
