/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/alex-guoba/tbd/internal/models"
	"github.com/alex-guoba/tbd/internal/provider"
	"github.com/alex-guoba/tbd/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Read prompt from pipe(like other file or shell cmd outputs), used to statics and anaylyze data
func readFromPipe() string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		logger.Error(err)
		return ""
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		// not pipe
		return ""
	} else {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			logger.Error(err)
			return ""
		}
		str := string(stdin)

		return strings.TrimSuffix(str, "\n")
	}
}

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat with ",
	Long: `A AI client for Baidu Qianfan. For example:
tbd "how about the weather?"
`,
	Run: func(cmd *cobra.Command, args []string) {
		pipeMsg := readFromPipe()

		if len(args) == 0 && pipeMsg == "" {
			logger.Warnf("Please input your prompt!")
			return
		}

		// Converge the chat completion request from pipe and args
		msg := args[0]
		if len(pipeMsg) > 0 {
			msg = pipeMsg + "\n" + msg
		}

		logger.Debug(msg)

		prov := provider.NewProvider()
		model := models.NewModel(viper.GetString("chat.model"), prov.GetClient())

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
