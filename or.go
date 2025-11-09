package main

import (
	"encoding/json"
	"io"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Provider struct {
	Only []string `json:"only"`
}

type Context struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Provider Provider  `json:"provider"`
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

type Body struct {
	pos     int
	Context *Context
}

func NewOR(model string, provider Provider) *Body {
	msgs := []Message{}

	ctx := &Context{
		Model:    model,
		Messages: msgs,
		Provider: provider,
	}

	return &Body{Context: ctx}
}

func (c *Context) NewUserMessage(message string) {
	c.Messages = append(c.Messages, Message{Role: "user", Content: message})
}

func (b *Body) Read(buf []byte) (int, error) {
	str, err := json.Marshal(b.Context)
	if err != nil {
		return 0, err
	}
	if b.pos >= len(str) {
		return 0, io.EOF
	}

	n := copy(buf, str[b.pos:])
	b.pos += n
	return n, nil
}
