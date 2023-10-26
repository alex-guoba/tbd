package provider

import (
	"context"

	ernie "github.com/anhao/go-ernie"
	"github.com/spf13/viper"
)

// / request
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Messages        []ChatCompletionMessage `json:"messages"`
	Temperature     float32                 `json:"temperature,omitempty"`
	TopP            float32                 `json:"top_p,omitempty"`
	PresencePenalty float32                 `json:"presence_penalty,omitempty"`
	Stream          bool                    `json:"stream"`
	UserId          string                  `json:"user_id,omitempty"`
	PenaltyScore    float32                 `json:"penalty_score,omitempty"`
	// Functions       []ErnieFunction         `json:"functions,omitempty"`
}

// resonse
// type ErniePluginUsage struct {
// 	Name           string `json:"name"`
// 	ParseTokens    int    `json:"parse_tokens"`
// 	AbstractTokens int    `json:"abstract_tokens"`
// 	SearchTokens   int    `json:"search_tokens"`
// 	TotalTokens    int    `json:"total_tokens"`
// }

// type ErnieUsage struct {
// 	PromptTokens     int                `json:"prompt_tokens"`
// 	CompletionTokens int                `json:"completion_tokens"`
// 	TotalTokens      int                `json:"total_tokens"`
// 	Plugins          []ErniePluginUsage `json:"plugins"`
// }

type ChatCompletionResponse struct {
	Id               string `json:"id"`
	Object           string `json:"object"`
	Created          int    `json:"created"`
	SentenceId       int    `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	// Usage            ErnieUsage `json:"usage"`
	// FunctionCall     ErnieFunctionCall `json:"function_call"`
	BanRound int `json:"ban_round"`
	APIError
}

type Model interface {
	ModelName(ctx context.Context) string
	GetCompletion(ctx context.Context, request *ChatCompletionRequest) (*ChatCompletionResponse, error)
}

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

func (prov *Provider) CreateChatCompletion(ctx context.Context, model Model, msg string) (*ChatCompletionResponse, error) {
	request := &ChatCompletionRequest{
		Messages: []ChatCompletionMessage{
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
