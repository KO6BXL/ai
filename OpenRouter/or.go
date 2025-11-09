package openrouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/KO6BXL/ai"
)

type DataCollection int

const (
	Allow DataCollection = iota
	Deny
)

const OR_CompletionsURL = "https://openrouter.ai/api/v1/chat/completions"

type Provider struct {
	Order                  []string       `json:"order"`
	AllowFallbacks         bool           `json:"allow_fallbacks"`
	RequireParameters      bool           `json:"require_parameters"`
	DataCollection         DataCollection `json:"data_collection"`
	Zdr                    bool           `json:"zdr"`
	EnforceDistillableText bool           `json:"enforce_distillable_text"`
	Only                   []string       `json:"only"`
	Ignore                 []string       `json:"ignore"`
	Quantizations          []string       `json:"quatizations"`
	Sort                   string         `json:"sort"`
}

type or_Mesg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type or_Request struct {
	Messages []or_Mesg `json:"messages"`
	Model    string    `json:"model"`
}

type OpenRouter struct {
	Messages []ai.Message
	Model    string
	Key      string
	Provider Provider
}

type ORResponse struct {
	Id      string  `json:"id"`
	Object  string  `json:"objext"`
	Created float64 `json:"created"`
	Model   string  `json:"model"`
	Usage   struct {
		InputTokens         float64 `json:"input_tokens"`
		OutputTokens        float64 `json:"output_tokens"`
		TotalTokens         float64 `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens float64 `json:"cached_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`

	Choices []struct {
		FinishReason       string  `json:"finish_reason"`
		NativeFinishReason string  `json:"native_finish_reason"`
		Index              float64 `json:"index"`
		Message            struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func NewOR(model string, apiKey string) *OpenRouter {
	msgs := []ai.Message{}

	ctx := &OpenRouter{
		Messages: msgs,
		Model:    model,
		Key:      apiKey,
	}

	return ctx
}

func (or *OpenRouter) SetProviders(provder Provider) {
	or.Provider = provder
}

func (or *OpenRouter) Request(ctx ai.Context) (ai.Response, error) {
	emptyresp := ai.Response{}
	client := http.Client{}
	if or.Key == "" {
		return emptyresp, errors.New("No OpenRouter API key")
	}
	if or.Model == "" {
		return emptyresp, errors.New("No Model Chosen")
	}
	msg := []or_Mesg{}

	for _, v := range ctx.Messages {
		msg = append(msg, or_Mesg{
			Role:    v.Role,
			Content: v.Content,
		})
	}

	reqO := or_Request{
		Messages: msg,
		Model:    or.Model,
	}
	reqB, err := json.Marshal(reqO)
	if err != nil {
		return emptyresp, err
	}
	body := strings.NewReader(fmt.Sprintf("%s\n", reqB))

	req, err := http.NewRequest(http.MethodPost, OR_CompletionsURL, body)
	if err != nil {
		return emptyresp, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", or.Key))
	req.Header.Set("Content-Type", "application/json")

	HTTPresp, err := client.Do(req)
	if err != nil {
		return emptyresp, err
	}
	StrResp, err := io.ReadAll(HTTPresp.Body)
	if err != nil {
		return emptyresp, err
	}

	if HTTPresp.StatusCode != http.StatusOK {
		return emptyresp, errors.New(fmt.Sprintf("Http error %d: %s\n", HTTPresp.StatusCode, StrResp))
	}
	orResp := &ORResponse{}
	err = json.Unmarshal(StrResp, orResp)
	if err != nil {
		return emptyresp, err
	}
	outs := []ai.Output{}

	for _, v := range orResp.Choices {
		outs = append(outs, ai.Output{
			FinishReason: v.FinishReason,
			Index:        int(v.Index),
			Message: ai.Message{
				Role:    v.Message.Role,
				Content: v.Message.Content,
			},
		})
	}

	return ai.Response{
		Id:    orResp.Id,
		Model: orResp.Model,
		Usage: ai.Usage{
			PromptTokens:     int(orResp.Usage.InputTokens),
			CompletionTokens: int(orResp.Usage.OutputTokens),
			TotalTokens:      int(orResp.Usage.TotalTokens),
		},
		Outputs: outs,
	}, nil
}
