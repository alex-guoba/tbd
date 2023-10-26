/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/alex-guoba/tbd/internal/models"
	"github.com/alex-guoba/tbd/internal/provider"
	"github.com/alex-guoba/tbd/pkg/logger"
	"github.com/spf13/cobra"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat with ",
	Long: `A AI client for Baidu Qianfan. For example:
tbd "how about the weather?"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Warnf("Please input your prompt")
			return
		}

		msg := args[0]
		prov := provider.NewProvider()
		model := models.NewErnieModel(prov.GetClient())

		// fmt.Println(msg)

		completion, err := prov.CreateChatCompletion(context.Background(), model, msg)
		if err != nil {
			logger.Errorf("ernie bot error: %v", err)
			return
		}
		logger.ChatReplay(completion.Result)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
