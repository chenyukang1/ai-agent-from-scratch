package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat/openai"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		fmt.Println("OPENAI_BASE_URL is not set")
		return
	}

	config := openai.ChatModelConfig{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		Model:      "deepseek/deepseek-v4-flash",
		HTTPClient: http.DefaultClient,
	}

	chatmodel := openai.NewChatModel(&config)
	input := []*schema.Message{
		{
			Role:    "user",
			Content: "What is the capital of France?",
		},
	}
	resp, err := chatmodel.Generate(context.Background(), input)
	if err != nil {
		fmt.Printf("Error generating response: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Content)
}
