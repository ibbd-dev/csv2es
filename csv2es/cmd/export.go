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
	"os"
	"strings"

	"github.com/ibbd-dev/go-csv"
	"github.com/spf13/cobra"
	"gopkg.in/olivere/elastic.v5"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export data to csv from es",
	Long: `export data to csv from es",
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
		conn, err := getESConnect()
		ctx := context.Background()
		exists, err := conn.IndexExists(cParams.IndexName).Do(ctx)
		if err != nil {
			panic(fmt.Errorf("check index exists error: %v", err.Error()))
		}

		if !exists {
			panic(fmt.Errorf("index %s is not exists", cParams.IndexName))
		}
		if cParams.Size <= 0 {
			cParams.Size = 1000
		}

		search := conn.Search(cParams.IndexName)
		if len(cParams.QueryField) > 0 {
			var query = elastic.NewTermQuery(cParams.QueryField, cParams.QueryValue)
			search.Query(query)
		}

		searchResult, err := search.Size(cParams.Size).Do(ctx)
		if err != nil {
			panic(fmt.Errorf("search error: %v", err.Error()))
		}

		var count int
		for i, hit := range searchResult.Hits.Hits {
			var row = make(map[string]interface{})
			err = json.Unmarshal(*hit.Source, &row)
			if err != nil {
				panic(fmt.Errorf("search %d: json unmarshal error: %v", i, err.Error()))
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

			count += 1
			var strRow = make(map[string]string)
			for k, v := range row {
				if sv, err := json.Marshal(v); err != nil {
					panic(fmt.Errorf("csv writer header error: %s", err.Error()))
				} else {
					strRow[k] = string(sv)
					strRow[k] = strings.Trim(strRow[k], "\"")
				}
			}
			writer.WriteRow(strRow)
		}
		writer.Flush()

		fmt.Printf("Total %d\n", count)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")
	exportCmd.PersistentFlags().StringVar(&cParams.QueryField, "query-field", "", "过滤字段")
	exportCmd.PersistentFlags().StringVar(&cParams.QueryValue, "query-value", "", "过滤字段对应的值")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// exportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
