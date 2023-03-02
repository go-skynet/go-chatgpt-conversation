package main

import (
	"bufio"
	"fmt"
	"os"

	conv "github.com/mudler/go-chatgpt-conversation"
)

func main() {
	conv, err := conv.NewConversation(os.Getenv("OPENAI_API_TOKEN"))
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Prompt:")

	for scanner.Scan() {
		line := scanner.Text()

		// do something with the line, for example, print it out
		data, err := conv.User(line)
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
		fmt.Println("Prompt:")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
