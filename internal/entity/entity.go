package entity

import "fmt"

const (
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

type APIError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	ID        string `json:"id"`
}

type RequestError struct {
	HTTPStatusCode int
	Err            error
}

func (e *APIError) Error() string {
	return e.ErrorMsg
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("error, status code: %d, message: %s", e.HTTPStatusCode, e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

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
