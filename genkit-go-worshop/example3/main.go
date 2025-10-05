package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"
)

// GreetingInput represents input for greeting generation
type GreetingInput struct {
	Name     string `json:"name" jsonschema:"description=The person's name"`
	Language string `json:"language" jsonschema:"description=Language for greeting (english, spanish, french)"`
}

// GreetingOutput represents the generated greeting
type GreetingOutput struct {
	Greeting string `json:"greeting" jsonschema:"description=The generated greeting"`
}

// JokeInput represents input for joke generation
type JokeInput struct {
	Topic string `json:"topic" jsonschema:"description=The topic for the joke"`
}

// JokeOutput represents the generated joke
type JokeOutput struct {
	Joke string `json:"joke" jsonschema:"description=The generated joke"`
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

	// Define a greeting flow
	greetingFlow := genkit.DefineFlow(g, "greeting", func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
		prompt := fmt.Sprintf("Create a friendly greeting for %s in %s. Keep it warm and welcoming.", input.Name, input.Language)

		response, err := genkit.Generate(ctx, g,
			ai.WithPrompt(prompt),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate greeting: %w", err)
		}

		return &GreetingOutput{
			Greeting: response.Text(),
		}, nil
	})

	// Define a joke generator flow
	jokeFlow := genkit.DefineFlow(g, "jokeGenerator", func(ctx context.Context, input *JokeInput) (*JokeOutput, error) {
		prompt := fmt.Sprintf("Create a clean, family-friendly joke about %s. Keep it short and funny.", input.Topic)

		response, err := genkit.Generate(ctx, g,
			ai.WithPrompt(prompt),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to generate joke: %w", err)
		}

		return &JokeOutput{
			Joke: response.Text(),
		}, nil
	})

	// Test the flows locally
	fmt.Println("=== Testing Greeting Flow ===")
	greetingResult, err := greetingFlow.Run(ctx, &GreetingInput{
		Name:     "Alice",
		Language: "english",
	})
	if err != nil {
		log.Printf("Error running greeting flow: %v", err)
	} else {
		fmt.Printf("Greeting Result: %s\n", greetingResult.Greeting)
	}

	fmt.Println("\n=== Testing Joke Flow ===")
	jokeResult, err := jokeFlow.Run(ctx, &JokeInput{
		Topic: "programming",
	})
	if err != nil {
		log.Printf("Error running joke flow: %v", err)
	} else {
		fmt.Printf("Joke Result: %s\n", jokeResult.Joke)
	}

	// Test multiple greetings
	fmt.Println("\n=== Testing Multiple Greetings ===")
	people := []GreetingInput{
		{Name: "Bob", Language: "spanish"},
		{Name: "Claire", Language: "french"},
	}

	for _, person := range people {
		result, err := greetingFlow.Run(ctx, &person)
		if err != nil {
			log.Printf("Error greeting %s: %v", person.Name, err)
			continue
		}
		fmt.Printf("%s (%s): %s\n", person.Name, person.Language, result.Greeting)
	}

	// Set up HTTP server to serve the flows
	mux := http.NewServeMux()
	mux.HandleFunc("POST /greeting", genkit.Handler(greetingFlow))
	mux.HandleFunc("POST /jokeGenerator", genkit.Handler(jokeFlow))

	// Print sample usage
	fmt.Println("\n=== Server Starting ===")
	fmt.Println("Starting server on http://localhost:3400")
	fmt.Println("Flows available at:")
	fmt.Println("  POST http://localhost:3400/greeting")
	fmt.Println("  POST http://localhost:3400/jokeGenerator")
	fmt.Println("\nSample curl commands:")
	fmt.Println(`  curl -X POST "http://localhost:3400/greeting" \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"data": {"name": "Alice", "language": "english"}}'`)
	fmt.Println()
	fmt.Println(`  curl -X POST "http://localhost:3400/jokeGenerator" \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"data": {"topic": "programming"}}'`)

	// Start the server
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}
