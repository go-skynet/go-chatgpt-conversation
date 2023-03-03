# go-chatgpt-conversation

A simple wrapper around the OpenAI api to provide a conversational model that keeps history.

For example:

```golang
package main

import (
	"bufio"
	"fmt"
	"os"

	conversation "github.com/mudler/go-chatgpt-conversation"
)

func main() {
	conv, err := conversation.New(
		os.Getenv("OPENAI_API_TOKEN"),
		conversation.WithInitialPrompt("You are a cat. You can reply with 'Meow' for yes, and 'Meow Meow' for no."),
	)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Prompt:")

	for scanner.Scan() {
		line := scanner.Text()

		// do something with the line, for example, print it out
		data, err := conv.Chat(line)
		if err != nil {
			panic(err)
		}

		// Save the conversation
		err = conv.Save("/tmp/conversation.json")

		// Load it back
		conv, err = conversation.Load("/tmp/conversation.json", os.Getenv("OPENAI_API_TOKEN"))
		fmt.Println(data)
		fmt.Println("Prompt:")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```