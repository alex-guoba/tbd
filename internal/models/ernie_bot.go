package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alex-guoba/tbd/internal/provider"
	ernie "github.com/anhao/go-ernie"
)

const ModelNameErnieBot = "ernie-bot"

type ErnieModel struct {
	client *ernie.Client
}

func NewErnieModel(client *ernie.Client) ErnieModel {
	return ErnieModel{
		client: client,
	}
}

func (m ErnieModel) ModelName(ctx context.Context) string {
	return ModelNameErnieBot
}

func (m ErnieModel) convertCompletionReplay(ersp *ernie.ErnieBotResponse) (*provider.ChatCompletionResponse, error) {

	buf, err := json.Marshal(ersp)
	if err != nil {
		return nil, fmt.Errorf("json error. %v", err.Error())
	}

	var completionRsp provider.ChatCompletionResponse
	err = json.Unmarshal(buf, &completionRsp)
	if err != nil {
		return nil, fmt.Errorf("json error. %v", err.Error())
	}

	return &completionRsp, nil
}

func (m ErnieModel) GetCompletion(ctx context.Context, request *provider.ChatCompletionRequest) (*provider.ChatCompletionResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	req := ernie.ErnieBotRequest{}
	req.Messages = []ernie.ChatCompletionMessage{
		{
			Role:    ernie.MessageRoleUser,
			Content: request.Messages[0].Content,
		},
	}

	rsp, err := m.client.CreateErnieBotChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return m.convertCompletionReplay(&rsp)
}
