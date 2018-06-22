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
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/olivere/elastic.v5"
)

type CommonParams struct {
	Debug bool
	Size  int // 导入导出的数量，0表示不限

	// es config
	Host      string
	Port      int
	IndexName string
	DocType   string

	// import
	DeleteIndex bool   // 是否删除原索引
	Mapping     string // mapping 文件

	// export
	QueryField string
	QueryValue string

	CsvFilename string
}

var cParams CommonParams

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "csv2es",
	Version: "v1.1",
	Short:   "import/export data beteen csv and es",
	Long: `import/export data between csv and es

实现功能：

- [x] import: 从csv文件导入数据到es
- [x] export: 从es导出数据到csv文件

Author:  Alex Cai
BuildAt: 20180622
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVar(&cParams.Debug, "debug", false, "debug mode")
	rootCmd.PersistentFlags().StringVar(&cParams.Host, "host", "localhost", "es host")
	rootCmd.PersistentFlags().IntVar(&cParams.Port, "port", 9200, "es port")
	rootCmd.PersistentFlags().StringVar(&cParams.IndexName, "index", "", "es index name")
	rootCmd.PersistentFlags().StringVar(&cParams.DocType, "type", "", "es doc type")
	rootCmd.PersistentFlags().StringVar(&cParams.CsvFilename, "csv", "", "csv filename")

	rootCmd.PersistentFlags().IntVar(&cParams.Size, "size", 0, "导入导出的数量，默认为0，表示导入不限量，导出限制为1000")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func getESConnect() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("http://%s:%d", cParams.Host, cParams.Port)),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5),
	)
}
