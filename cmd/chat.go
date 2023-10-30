/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/internal/models"
	"github.com/alex-guoba/tbd/internal/prompts"
	"github.com/alex-guoba/tbd/internal/provider"
	"github.com/alex-guoba/tbd/pkg/logger"
	"github.com/alex-guoba/tbd/pkg/sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagInteract bool

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "chat with ",
	Long: `AI client for Baidu Qianfan. For example:
tbd "how about the weather?"
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		promp := prompts.New(ctx)

		mood, piped := promp.ReadFromPipe()
		if mood && flagInteract {
			logger.Error("Interactive mode can't be used with piping.")
			cancel()
			return
		}

		agent := provider.NewAgent()
		prov := provider.NewProvider()
		model := models.NewModel(viper.GetString("chat.model"), prov.GetClient())

		msg := ""
		if len(args) > 0 {
			msg = args[0]
		}
		msg = promp.Converge(msg, piped)

		// sync subscribe
		nsync := notionSync(ctx, agent)
		memos := memosSync(ctx, agent)

		// Trigger graceful shutdown on SIGINT or SIGTERM.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-c
			logger.Infof("%s recevied. exist", sig.String())
			cancel()
			agent.Close()
			closeSyncAll(nsync, memos)
			// stdin blocking, can't be stoped —— todo: fix it
			os.Exit(0)
		}()

		startDialogs(ctx, prov, model, msg, flagInteract, agent)

		// stop every thing
		agent.Close()
		closeSyncAll(nsync, memos)
	},
}

func closeSyncAll(nsync *sync.SyncNotion, memos *sync.SyncMemos) {
	if nsync != nil {
		nsync.Close()
	}
	if memos != nil {
		memos.Close()
	}
}

// synchornize chat message to notion
func notionSync(ctx context.Context, agent *provider.Agent) *sync.SyncNotion {
	token := viper.GetString("sync.notion.token")
	db := viper.GetString("sync.notion.token")

	if len(token) == 0 || len(db) == 0 {
		return nil
	}

	logger.Debug("notion subscribed")

	ch := agent.Subscribe(provider.TopicChat)
	notion := sync.NewNotion(ctx, token, db)

	notion.Consume(ch)
	return notion
}

// synchornize chat message to memos
func memosSync(ctx context.Context, agent *provider.Agent) *sync.SyncMemos {
	host := viper.GetString("sync.memos.host")
	token := viper.GetString("sync.memos.token")

	if len(host) == 0 || len(token) == 0 {
		return nil
	}

	logger.Debug("memos subscribed")

	ch := agent.Subscribe(provider.TopicChat)
	memos := sync.NewMemos(ctx, token, host)

	memos.Consume(ch)
	return memos
}

// Start dialogus(once or mutlti-round)
func startDialogs(ctx context.Context, prov *provider.Provider, model models.Model,
	initialMsg string, interact bool, agent *provider.Agent) {
	request := &entity.ChatCompletionRequest{}
	msg := initialMsg
	promp := prompts.New(ctx)

	if len(msg) == 0 {
		if !interact {
			logger.Error("Please input prompt meesage!")
			return
		}
		msg = promp.ReadFromStdin()
	}

	for {
		if promp.ShouldStop(msg) {
			logger.Info("Bye.")
			break
		}

		request.Messages = append(request.Messages, entity.ChatCompletionMessage{
			Role:    entity.MessageRoleUser,
			Content: msg,
		})

		logger.Debug(request.Messages)

		completion, err := prov.CreateChatCompletion(ctx, model, request, agent)
		if err != nil {
			logger.Errorf("ernie bot error: %v", err)
			return
		}

		logger.ChatReplay(completion.Result)

		if !interact {
			break
		}
		request.Messages = append(request.Messages, entity.ChatCompletionMessage{
			Role:    entity.MessageRoleAssistant,
			Content: completion.Result,
		})
		msg = promp.ReadFromStdin()
	}
}

func init() {
	// May be promote to root cmd.
	chatCmd.Flags().BoolVarP(&flagInteract, "interact", "i", false, "Interactive mode for multi-round dialogus")

	rootCmd.AddCommand(chatCmd)
}
