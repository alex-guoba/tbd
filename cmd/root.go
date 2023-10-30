/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tbd",
	Short: "A AI client for Baidu Qianfan.",
	Long: `A AI client for Baidu Qianfan. For example:
tbd "how about the weather?"
`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(defCmd string) {
	// var cmdFound bool
	// cmd := rootCmd.Commands()

	// for _, a := range cmd {
	// 	for _, b := range os.Args[1:] {
	// 		if a.Name() == b {
	// 			cmdFound = true
	// 			break
	// 		}
	// 	}
	// }
	// if !cmdFound {
	// 	args := append([]string{defCmd}, os.Args[1:]...)
	// 	rootCmd.SetArgs(args)
	// }
	if err := rootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tbd.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
