package ernie_bot_4

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/pkg/logger"
	ernie "github.com/anhao/go-ernie"
	// "github.com/lxc/lxd/shared/logger"
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

// merge stream replay
// replay format: https://cloud.baidu.com/doc/WENXINWORKSHOP/s/clntwmv7t#%E5%93%8D%E5%BA%94%E7%A4%BA%E4%BE%8B%EF%BC%88%E6%B5%81%E5%BC%8F%EF%BC%89
func (m ErnieBot4) mergeStreamReplay(last *entity.ChatCompletionResponse, current *entity.ChatCompletionResponse) *entity.ChatCompletionResponse {
	ret := &entity.ChatCompletionResponse{}
	*ret = *current // copy all filed
	if last != nil {
		logger.Debug("[current replay]", current.Result)
		ret.Result = last.Result + current.Result // except Result
	}
	return ret
}

func (m ErnieBot4) GetCompletion(ctx context.Context, request *entity.ChatCompletionRequest) (*entity.ChatCompletionResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	req := ernie.ErnieBot4Request{
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
	// req.Messages = []ernie.ChatCompletionMessage{
	// 	{
	// 		Role:    ernie.MessageRoleUser,
	// 		Content: request.Messages[0].Content,
	// 	},
	// }

	rsp, err := m.client.CreateErnieBot4ChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	return m.convertCompletionReplay(&rsp)
}

func (m ErnieBot4) GetCompletionStream(ctx context.Context, request *entity.ChatCompletionRequest,
	callback entity.StreamCallback) (*entity.ChatCompletionResponse, error) {
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("empty message")
	}
	var current, last *entity.ChatCompletionResponse

	req := ernie.ErnieBot4Request{}
	for _, msg := range request.Messages {
		req.Messages = append(req.Messages, ernie.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	stream, err := m.client.CreateErnieBot4ChatCompletionStream(ctx, req)
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

		current, err = m.convertCompletionReplay(&rsp)
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
