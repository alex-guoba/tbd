package provider

import (
	"context"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/internal/models"
	"github.com/alex-guoba/tbd/pkg/animator"
	"github.com/alex-guoba/tbd/pkg/logger"
	ernie "github.com/anhao/go-ernie"
	"github.com/spf13/viper"
)

// Chat message
const TopicChat = "chat"

type TopicChatMsg struct {
	Req *entity.ChatCompletionRequest
	Rsp *entity.ChatCompletionResponse
}

///////////////////////////////////////////////////////////////////

type Provider struct {
	client *ernie.Client
}

func NewProvider() *Provider {
	client := ernie.NewDefaultClient(viper.GetString("appkey"), viper.GetString("secretkey"))
	logger.Debugf("key: %s, secret: %s", viper.GetString("appkey"), viper.GetString("secretkey"))

	return &Provider{client: client}
}

func (prov *Provider) GetClient() *ernie.Client {
	return prov.client
}

func (prov *Provider) CreateChatCompletion(ctx context.Context, model models.Model,
	request *entity.ChatCompletionRequest, agent *Agent) (*entity.ChatCompletionResponse, error) {
	spinner := animator.New()

	spinner.Start()
	defer spinner.Stop()

	if rsp, err := model.GetCompletion(ctx, request); err != nil {
		return nil, err
	} else {
		// publish
		if agent != nil {
			agent.Publish(TopicChat, &TopicChatMsg{
				Req: request,
				Rsp: rsp,
			})
		}

		return rsp, nil
	}
}
