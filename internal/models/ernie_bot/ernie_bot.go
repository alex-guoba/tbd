package ernie_bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/pkg/logger"
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

// ERNIE-Bot, ERNIE-Bot-turbo share the replay struct, should fix it for go-ernie
func (m ErnieBot) convertCompletionTurboReplay(ersp *ernie.ErnieBotTurboResponse) (*entity.ChatCompletionResponse, error) {
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

	// TODO: support `system` parameter
	req := ernie.ErnieBotRequest{}
	for _, msg := range request.Messages {
		req.Messages = append(req.Messages, ernie.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	rsp, err := m.client.CreateErnieBotChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return m.convertCompletionReplay(&rsp)
}

// merge stream replay
// replay format: https://cloud.baidu.com/doc/WENXINWORKSHOP/s/jlil56u11#%E8%AF%B7%E6%B1%82%E7%A4%BA%E4%BE%8B%EF%BC%88%E6%B5%81%E5%BC%8F%EF%BC%89
func (m ErnieBot) mergeStreamReplay(last *entity.ChatCompletionResponse, current *entity.ChatCompletionResponse) *entity.ChatCompletionResponse {
	ret := &entity.ChatCompletionResponse{}
	*ret = *current // copy all filed
	if last != nil {
		ret.Result = last.Result + current.Result // except Result
	}
	return ret
}

func (m ErnieBot) GetCompletionStream(ctx context.Context, request *entity.ChatCompletionRequest,
	callback entity.StreamCallback) (*entity.ChatCompletionResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("empty message")
	}
	var current, last *entity.ChatCompletionResponse

	req := ernie.ErnieBotRequest{
		Temperature:  request.Temperature,
		TopP:         request.TopP,
		Stream:       request.Stream,
		PenaltyScore: request.PenaltyScore,
		UserId:       request.UserId,
	}
	for _, msg := range request.Messages {
		req.Messages = append(req.Messages, ernie.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	stream, err := m.client.CreateErnieBotChatCompletionStream(ctx, req)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	for {
		rsp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			logger.Debug("ernie bot Stream finished")
			break
		}
		if err != nil {
			return nil, err
		}

		current, err = m.convertCompletionTurboReplay(&rsp)
		if err != nil {
			return nil, err
		}

		if callback != nil {
			if err = callback(current); err != nil {
				// callback error
				return nil, err
			}
		}

		last = m.mergeStreamReplay(last, current)
	}

	return last, nil
}
