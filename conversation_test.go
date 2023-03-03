package chatgpt_test

import (
	"os"
	"path/filepath"

	. "github.com/mudler/go-chatgpt-conversation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/otiai10/openaigo"
)

var _ = Describe("ask to gpt", func() {
	Context("conversation", func() {
		It("can track history", func() {
			if os.Getenv("OPENAI_API_TOKEN") == "" {
				Skip("OPENAI_API_TOKEN required to run this test")
			}
			cat, err := NewConversation(os.Getenv("OPENAI_API_TOKEN"),
				WithInitialPrompt("You are a cat. You can reply with 'Meow' for yes, and 'Meow Meow' for no."))
			Expect(err).ToNot(HaveOccurred())

			res, err := cat.User("Are you a cat?")
			Expect(err).ToNot(HaveOccurred())

			Expect(res).To(Equal("Meow."))

			res, err = cat.User("Are you a human?")
			Expect(err).ToNot(HaveOccurred())

			Expect(res).To(Equal("Meow Meow."))
		})

		It("can save history", func() {
			testHistory := []openaigo.ChatMessage{
				{Role: "system", Content: "foo"},
				{Role: "user", Content: "bar"},
				{Role: "assistant", Content: "response"}}

			cat, err := NewConversation("none",
				WithHistory(
					testHistory,
				),
			)
			Expect(err).ToNot(HaveOccurred())
			dir, err := os.MkdirTemp("", "test")
			Expect(err).ToNot(HaveOccurred())
			defer os.RemoveAll(dir)

			f := filepath.Join(dir, "test.json")

			err = cat.Save(f)
			Expect(err).ToNot(HaveOccurred())

			conv, err := Load(f, "test")
			Expect(err).ToNot(HaveOccurred())

			Expect(conv.History).To(Equal(testHistory))
		})
	})
})
