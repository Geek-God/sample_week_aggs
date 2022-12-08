// Package sample_search_export
// @author:WXZ
// @date:2022/11/30
// @note

package main

import (
	"github.com/spf13/viper"
	"log"
	"sample_week_aggs/cmd"
	elasticInit "sample_week_aggs/initd/elasticInitd"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags)
	err := elasticInit.New()
	if err != nil {
		log.Fatal(err)
	}
	err = timeLoc()
	if err != nil {
		log.Fatal(err)
	}

}
func main() {
	cmd.Execute()
}

// timeLoc
// @Author WXZ
// @Description: //TODO 设置时区
func timeLoc() error {
	loc, err := time.LoadLocation(viper.GetString("timezone"))
	if err != nil {
		return err
	}
	time.Local = loc
	return nil
}
