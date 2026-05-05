package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
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

	config := openai.DefaultConfig(apiKey)
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		config.BaseURL = baseURL
	}

	config.HTTPClient = &http.Client{
		Transport: &headerTransport{
			rt: http.DefaultTransport,
		},
		Timeout: 10 * time.Second,
	}

	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "deepseek/deepseek-v4-flash",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

type headerTransport struct {
	rt http.RoundTripper
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// req.Header.Set("HTTP-Referer", "https://test.com")
	// req.Header.Set("X-OpenRouter-Title", "test")
	req.Header.Set("User-Agent", "curl/7.81.0")
	req.Header.Set("Accept", "*/*")
	// 打印 URL
	fmt.Printf("Request URL: %s %s\n", req.Method, req.URL.String())
	// 打印 Headers
	for k, v := range req.Header {
		fmt.Printf("Header: %s = %v\n", k, v)
	}
	// 打印 Body
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 写回 Body
		fmt.Printf("Body: %s\n", string(bodyBytes))
	}

	return h.rt.RoundTrip(req)
}
