package chatgpt

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/otiai10/openaigo"
)

const (
	roleSystem    = "system"
	roleUser      = "user"
	roleAssistant = "assistant"

	defaultModel = "gpt-3.5-turbo"
)

func (c *Conversation) converse(role, input string) (string, error) {
	c.History = append(c.History, openaigo.ChatMessage{Role: role, Content: input})
	resp, err := chatGPT(context.Background(), c.model, c.apiKey, c.History)
	if err != nil {
		return "", err
	}
	c.History = append(c.History, openaigo.ChatMessage{Role: roleAssistant, Content: resp})
	return resp, err
}

// System writes in the conversation as the "system" role.
// You can use a system level instruction to guide your model's behavior throughout the conversation.
// see: https://help.openai.com/en/articles/7042661-chatgpt-api-transition-guide
func (c *Conversation) System(input string) (string, error) {
	return c.converse(roleSystem, input)
}

// User writes in the conversation as the "user" role.
func (c *Conversation) User(input string) (string, error) {
	return c.converse(roleUser, input)
}

func chatGPT(ctx context.Context, model string, apiKey string, history []openaigo.ChatMessage) (string, error) {
	client := openaigo.NewClient(apiKey)
	request := openaigo.ChatCompletionRequestBody{
		Model:    model,
		Messages: history,
	}

	resp, err := client.Chat(ctx, request)
	if err != nil {
		return "", err

	}

	if len(resp.Choices) != 1 {
		return "", fmt.Errorf("invalid response: %+v", resp)
	}
	return resp.Choices[0].Message.Content, nil
}

// Conversation is the abstract conversation view.
type Conversation struct {
	History []openaigo.ChatMessage
	apiKey  string
	model   string
}

// Option is the options for the conversation
type Option func(c *Conversation) error

// WithInitialPrompt allow to load an initial prompt to guide your model's behavior throughout the conversation.
// The prompt will be written as "system" role.
func WithInitialPrompt(s string) Option {
	return func(c *Conversation) error {
		c.History = append(c.History, openaigo.ChatMessage{Role: roleSystem, Content: s})
		return nil
	}
}

// WithModel allows to override the default model (gpt-3.5-turbo)
func WithModel(model string) Option {
	return func(c *Conversation) error {
		c.model = model
		return nil
	}
}

// WithHistory allows to inject an history in the conversation
func WithHistory(msg []openaigo.ChatMessage) Option {
	return func(c *Conversation) error {
		c.History = msg
		return nil
	}
}

// NewConversation starts a new conversation.
// Requires an OpenAI API key, and an optional list of options.
// It returns a new conversation and an error if there is a failure.
func NewConversation(apiKey string, opts ...Option) (*Conversation, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("an API key is required")
	}
	c := &Conversation{apiKey: apiKey, model: defaultModel}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Save the conversation to disk at the given full path.
func (c *Conversation) Save(dst string) error {
	dat, err := json.Marshal(c.History)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, dat, 0655)
}

// Load loads a conversation from the disk with the given options.
func Load(file, apiKey string, opts ...Option) (*Conversation, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	msgs := []openaigo.ChatMessage{}

	if err := json.Unmarshal(b, &msgs); err != nil {
		return nil, err
	}

	return NewConversation(apiKey, append([]Option{WithHistory(msgs)}, opts...)...)
}

// var _ io.ReadWriter = Conversation{}

// func (c *Conversation) Read(p []byte) (n int, err error) {

// }
// Write(p []byte) (n int, err error)
