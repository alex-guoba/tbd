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
		if agent != nil {
			// deep copy request and response
			agent.Publish(TopicChat, &TopicChatMsg{
				Req: entity.CopyChatCompletionRequest(request),
				Rsp: entity.CopyChatCompletionResponse(rsp),
			})
		}

		return rsp, nil
	}
}

func (prov *Provider) CreateChatStreamCompletion(ctx context.Context, model models.Model,
	request *entity.ChatCompletionRequest, agent *Agent, callback entity.StreamCallback) (*entity.ChatCompletionResponse, error) {

	if rsp, err := model.GetCompletionStream(ctx, request, callback); err != nil {
		return nil, err
	} else {
		if agent != nil {
			// deep copy request and response
			agent.Publish(TopicChat, &TopicChatMsg{
				Req: entity.CopyChatCompletionRequest(request),
				Rsp: entity.CopyChatCompletionResponse(rsp),
			})
		}

		return rsp, nil
	}
}
