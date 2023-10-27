package provider

import (
	"context"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/internal/models"
	ernie "github.com/anhao/go-ernie"
	"github.com/spf13/viper"
)

/////////////////////////////////////////////////////////////////////

type Provider struct {
	client *ernie.Client
}

func NewProvider() *Provider {
	client := ernie.NewDefaultClient(viper.GetString("appkey"), viper.GetString("appsecret"))
	return &Provider{client: client}
}

func (prov *Provider) GetClient() *ernie.Client {
	return prov.client
}

func (prov *Provider) CreateChatCompletion(ctx context.Context, model models.Model, msg string) (*entity.ChatCompletionResponse, error) {
	request := &entity.ChatCompletionRequest{
		Messages: []entity.ChatCompletionMessage{
			{
				Content: msg,
			},
		},
		// ignore other field now
	}

	if rsp, err := model.GetCompletion(ctx, request); err != nil {
		return nil, err
	} else {
		return rsp, nil
	}
}
