package conversation

import (
	"context"

	"github.com/otiai10/openaigo"
)

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

// WithContext associate a context to the conversation
func WithContext(ctx context.Context) Option {
	return func(c *Conversation) error {
		c.ctx = ctx
		return nil
	}
}
