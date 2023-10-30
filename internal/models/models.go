package models

import (
	"context"
	"strings"

	"github.com/alex-guoba/tbd/internal/entity"
	"github.com/alex-guoba/tbd/internal/models/ernie_bot"
	"github.com/alex-guoba/tbd/internal/models/ernie_bot_4"
	"github.com/alex-guoba/tbd/pkg/logger"

	ernie "github.com/anhao/go-ernie"
)

type Model interface {
	Name(ctx context.Context) string

	// single replay
	GetCompletion(ctx context.Context, request *entity.ChatCompletionRequest) (*entity.ChatCompletionResponse, error)

	// stream replay
	GetCompletionStream(ctx context.Context, request *entity.ChatCompletionRequest,
		callback entity.StreamCallback) (*entity.ChatCompletionResponse, error)
}

func NewModel(name string, client *ernie.Client) Model {
	lowerName := strings.ToLower(name)

	switch lowerName {
	case ernie_bot.ModelNameErnieBot:
		return ernie_bot.New(client)

	case ernie_bot_4.ModelNameErnieBot4:
		return ernie_bot_4.New(client)

	default:
		logger.Errorf("unknown model name: %s, use default", name)
		return ernie_bot_4.New(client)
	}
}
