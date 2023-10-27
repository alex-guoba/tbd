package ernie_bot

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alex-guoba/tbd/internal/entity"
	ernie "github.com/anhao/go-ernie"
)

const ModelNameErnieBot = "ernie-bot"

type ErnieBot struct {
	client *ernie.Client
}

func New(client *ernie.Client) ErnieBot {
	return ErnieBot{
		client: client,
	}
}

func (m ErnieBot) Name(ctx context.Context) string {
	return ModelNameErnieBot
}

func (m ErnieBot) convertCompletionReplay(ersp *ernie.ErnieBotResponse) (*entity.ChatCompletionResponse, error) {
	buf, err := json.Marshal(ersp)
	if err != nil {
		return nil, fmt.Errorf("json error. %v", err.Error())
	}

	var completionRsp entity.ChatCompletionResponse
	err = json.Unmarshal(buf, &completionRsp)
	if err != nil {
		return nil, fmt.Errorf("json error. %v", err.Error())
	}

	return &completionRsp, nil
}

func (m ErnieBot) GetCompletion(ctx context.Context, request *entity.ChatCompletionRequest) (*entity.ChatCompletionResponse, error) {
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
