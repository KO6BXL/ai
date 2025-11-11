package ai

type Context struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string
	Content string
}

type Output struct {
	FinishReason string
	Index        int
	Message      Message
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type Response struct {
	Id      string
	Model   string
	Outputs []Output
	Usage   Usage
}

type Driver interface {
	Request(ctx Context) (Response, error)
}

type AI struct {
	ctx    Context
	driver Driver
}

func NewAI(driver Driver) *AI {
	msgs := []Message{}
	ctx := Context{
		Messages: msgs,
	}
	return &AI{ctx, driver}
}

func (ai *AI) Prompt(message string) {
	ai.ctx.Messages = append(ai.ctx.Messages, Message{Role: "system", Content: message})
}

func (ai *AI) Message(message string) (Response, error) {
	ai.ctx.Messages = append(ai.ctx.Messages, Message{Role: "user", Content: message})

	resp, err := ai.driver.Request(ai.ctx)
	if err != nil {
		return resp, err
	}
	ai.ctx.Messages = append(ai.ctx.Messages, Message{Role: "assistant", Content: resp.Outputs[len(resp.Outputs)-1].Message.Content})
	return resp, err
}
