package main

import (
	"fmt"
	"log"
	"time"

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
	rootCmd.PersistentFlags().IntVarP(&Conc, "conc", "c", 1, "并发数")
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
		fmt.Println("Conc:", Conc)

		downloader.WithProgress(Verbose)

		start := time.Now()
		err := downloader.Download(Url, Out)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("耗时", time.Since(start))
	},
}
