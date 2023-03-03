package conversation

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

// Conversation is the abstract conversation view.
type Conversation struct {
	History []openaigo.ChatMessage
	apiKey  string
	model   string
	ctx     context.Context
}

// New creates a new conversation.
// Requires an OpenAI API key, and an optional list of options.
// It returns a new conversation and an error if there is a failure.
func New(apiKey string, opts ...Option) (*Conversation, error) {
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

	return New(apiKey, append([]Option{WithHistory(msgs)}, opts...)...)
}

// Save the conversation to disk at the given full path.
func (c *Conversation) Save(dst string) error {
	dat, err := json.Marshal(c.History)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, dat, 0655)
}

func (c *Conversation) apiCall(ctx context.Context, history []openaigo.ChatMessage) (string, error) {
	client := openaigo.NewClient(c.apiKey)
	request := openaigo.ChatCompletionRequestBody{
		Model:    c.model,
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

func (c *Conversation) converse(ctx context.Context, role, input string) (string, error) {
	c.History = append(c.History, openaigo.ChatMessage{Role: role, Content: input})
	resp, err := c.apiCall(ctx, c.History)
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
	return c.converse(c.ctx, roleSystem, input)
}

// User writes in the conversation as the "user" role.
func (c *Conversation) User(input string) (string, error) {
	return c.converse(c.ctx, roleUser, input)
}

// Chat is syntax sugar. is equivalent to call User
func (c *Conversation) Chat(input string) (string, error) {
	return c.User(input)
}
