/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/alex-guoba/tbd/cmd"
	"github.com/alex-guoba/tbd/config"
)

func main() {
	// init()
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	cmd.Execute("chat")
}
