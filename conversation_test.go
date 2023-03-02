package chatgpt_test

import (
	"os"

	. "github.com/mudler/go-chatgpt-conversation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ask to gpt", func() {
	Context("conversation", func() {
		It("can track history", func() {
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
	})
})
