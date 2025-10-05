package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func main() {
	ctx := context.Background()

	// Load environment variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Initialize Genkit with Google AI plugin
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"),
	)

	// Example 1: Generate a welcome message
	fmt.Println("=== Welcome Message ===")
	response1, err := genkit.Generate(ctx, g,
		ai.WithPrompt("Write a short welcome message for a new team member joining our development team."),
	)
	if err != nil {
		log.Printf("Error generating welcome message: %v", err)
	} else {
		fmt.Println(response1.Text())
	}

	// Example 2: Create a simple task list
	fmt.Println("\n=== Task List ===")
	response2, err := genkit.Generate(ctx, g,
		ai.WithPrompt("Create a simple 3-item todo list for setting up a new development environment."),
	)
	if err != nil {
		log.Printf("Error generating task list: %v", err)
	} else {
		fmt.Println(response2.Text())
	}

	// Example 3: Write a brief explanation
	fmt.Println("\n=== Explanation ===")
	response3, err := genkit.Generate(ctx, g,
		ai.WithPrompt("Explain what Go programming language is in 2-3 simple sentences."),
	)
	if err != nil {
		log.Printf("Error generating explanation: %v", err)
	} else {
		fmt.Println(response3.Text())
	}
}
