package conversation_test

import (
	. "github.com/mudler/go-chatgpt-conversation"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/otiai10/openaigo"
)

var _ = Describe("test options", func() {
	Context("constructor", func() {
		It("can set history", func() {
			testHistory := []openaigo.ChatMessage{
				{Role: "system", Content: "foo"},
				{Role: "user", Content: "bar"},
				{Role: "assistant", Content: "response"}}

			cat, err := New("none",
				WithHistory(
					testHistory,
				),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(cat.History).To(Equal(testHistory))
		})
		It("can set initial prompt", func() {
			expectedHistory := []openaigo.ChatMessage{
				{Role: "system", Content: "foo"},
			}
			cat, err := New("none",
				WithInitialPrompt("foo"),
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(cat.History).To(Equal(expectedHistory))
		})
	})
})
