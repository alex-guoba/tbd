package ernie_bot_4

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alex-guoba/tbd/internal/entity"
	ernie "github.com/anhao/go-ernie"
)

const ModelNameErnieBot4 = "ernie-bot-4"

type ErnieBot4 struct {
	client *ernie.Client
}

func New(client *ernie.Client) ErnieBot4 {
	return ErnieBot4{
		client: client,
	}
}

func (m ErnieBot4) Name(ctx context.Context) string {
	return ModelNameErnieBot4
}

func (m ErnieBot4) convertCompletionReplay(ersp *ernie.ErnieBot4Response) (*entity.ChatCompletionResponse, error) {
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

func (m ErnieBot4) GetCompletion(ctx context.Context, request *entity.ChatCompletionRequest) (*entity.ChatCompletionResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	req := ernie.ErnieBot4Request{}
	req.Messages = []ernie.ChatCompletionMessage{
		{
			Role:    ernie.MessageRoleUser,
			Content: request.Messages[0].Content,
		},
	}

	rsp, err := m.client.CreateErnieBot4ChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return m.convertCompletionReplay(&rsp)
}
