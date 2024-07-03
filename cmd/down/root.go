package main

import (
	"fmt"
	"log"

	"github.com/favbox/gopkg/util/downloader"
	"github.com/spf13/cobra"
)

var Verbose bool
var Conc int
var Url string
var Out string

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "详细输出")
	rootCmd.PersistentFlags().StringVarP(&Url, "url", "u", "", "资源网址")
	rootCmd.PersistentFlags().StringVarP(&Out, "out", "o", "", "输出路径")
	rootCmd.PersistentFlags().IntVarP(&Conc, "conc", "c", 0, "并发数")
	err := rootCmd.MarkPersistentFlagRequired("url")
	if err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "down",
	Short: "并发资源下载器",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Verbose:", Verbose)
		fmt.Println("Url:", Url)
		fmt.Println("Out:", Out)
		if Conc == 0 {
			fmt.Println("Conc: 自动")
		} else {
			fmt.Println("Conc:", Conc)
		}

		downloader.WithProgress(Verbose)

		err := downloader.DownloadWithChunks(Url, Conc, Out)
		if err != nil {
			log.Fatal(err)
		}
	},
}
