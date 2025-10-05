package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

// Story represents a generated story with structured data
type Story struct {
	Title      string   `json:"title" jsonschema:"description=The title of the story"`
	Characters []string `json:"characters" jsonschema:"description=Main characters in the story"`
	Setting    string   `json:"setting" jsonschema:"description=Where the story takes place"`
	Plot       string   `json:"plot" jsonschema:"description=The main plot of the story"`
	Mood       string   `json:"mood" jsonschema:"description=The overall mood (happy, adventurous, mysterious, etc.)"`
}

// StoryInput represents the input parameters for story generation
type StoryInput struct {
	Character string `json:"character" jsonschema:"description=Main character for the story"`
	Setting   string `json:"setting" jsonschema:"description=Setting where the story takes place"`
}

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

	// Generate a structured story
	fmt.Println("=== Generating Structured Story ===")

	input := StoryInput{
		Character: "a friendly robot",
		Setting:   "a magical library",
	}

	fmt.Printf("Character: %s\n", input.Character)
	fmt.Printf("Setting: %s\n", input.Setting)

	prompt := fmt.Sprintf(`Create a short, fun story about %s in %s. 
		Make it light-hearted and suitable for all ages.
		Return the response as structured data with title, characters, setting, plot, and mood.`,
		input.Character, input.Setting)

	// Generate structured data using Go structs
	story, _, err := genkit.GenerateData[Story](ctx, g,
		ai.WithPrompt(prompt),
	)
	if err != nil {
		log.Fatalf("Error generating story: %v", err)
	}

	// Print the structured story
	fmt.Println("\n=== Generated Story ===")
	storyJSON, _ := json.MarshalIndent(story, "", "  ")
	fmt.Println(string(storyJSON))

	// Generate different stories
	fmt.Println("\n=== Generating Different Stories ===")

	stories := []StoryInput{
		{Character: "a curious cat", Setting: "a space station"},
		{Character: "a young wizard", Setting: "a bustling marketplace"},
		{Character: "a brave explorer", Setting: "an underwater city"},
	}

	for i, storyInput := range stories {
		fmt.Printf("\n--- Story %d: %s in %s ---\n", i+1, storyInput.Character, storyInput.Setting)

		prompt := fmt.Sprintf(`Create a short, fun story about %s in %s. 
			Make it light-hearted and suitable for all ages.
			Return the response as structured data with title, characters, setting, plot, and mood.`,
			storyInput.Character, storyInput.Setting)

		story, _, err := genkit.GenerateData[Story](ctx, g,
			ai.WithPrompt(prompt),
		)
		if err != nil {
			log.Printf("Error generating story %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Title: %s\n", story.Title)
		fmt.Printf("Plot: %s\n", story.Plot)
		fmt.Printf("Mood: %s\n", story.Mood)
	}
}
