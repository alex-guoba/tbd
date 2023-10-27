/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/alex-guoba/tbd/cmd"
	"github.com/alex-guoba/tbd/config"
	"github.com/alex-guoba/tbd/pkg/logger"
	"github.com/spf13/viper"
)

func main() {
	// init()
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	logger.SetLevel(viper.GetInt("log.level"))

	cmd.Execute("chat")
}
