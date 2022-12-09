// Package cmd
// @author:WXZ
// @date:2021/9/23
// @note

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sample_week_aggs/service/SampleWeekAggs"
	"time"
)

var (
	begin_time string
	end_time   string
)
var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Usage: dl.exe <OPTIONS> [id/hash/file]...",
	Long:  `HuoRong sample statistics`,
	Run:   runCmd,
}

func init() {
	rootCmd.Flags().StringVarP(&begin_time, "begin", "b", "", "begin time eg:\"2006-01-02 15:04:05\"")
	rootCmd.Flags().StringVarP(&end_time, "end", "e", "", "end time eg:\"2006-01-02 15:04:05\"")
}

// Execute
// @Author WXZ
// @Description: //TODO
func Execute() {
	rootCmd.Execute()
}

// runCmd
// @Author WXZ
// @Description: //TODO
// @param cmd *cobra.Command
// @param args []string
func runCmd(cmd *cobra.Command, args []string) {
	_, err := time.Parse("2006-01-02 15:04:05", begin_time)
	if err != nil {
		log.Fatal("Please enter the correct begin time format")
	}
	_, err = time.Parse("2006-01-02 15:04:05", end_time)
	if err != nil {
		log.Fatal("Please enter the correct end time format")
	}

	aggs := SampleWeekAggs.New(SampleWeekAggs.WitchTime(begin_time, end_time))
	err = aggs.MakeWord()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("SUCCESS")
}
