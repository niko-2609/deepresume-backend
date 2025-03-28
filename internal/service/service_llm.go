package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// LLMService provides functionality to interact with LLM models
type LLMService struct {
	client  *http.Client
	baseURL string
}

// NewLLMService creates a new LLMService instance
func NewLLMService() *LLMService {
	return &LLMService{
		client: &http.Client{
			Timeout: 120 * time.Second, // Set a reasonable timeout for LLM requests
		},
		baseURL: "http://localhost:11434/api", // Local Deepseek/Ollama instance
	}
}

// LLMRequest represents a request to the LLM API
type LLMRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// LLMResponse represents a response from the LLM API (Ollama format)
type LLMResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// GenerateContent sends a prompt to the LLM model and returns the generated text
func (s *LLMService) GenerateContent(ctx context.Context, model string, prompt string) (string, error) {
	// Create request
	reqBody := LLMRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // Not streaming for this method
	}

	// Convert request to JSON
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/generate", bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.New("LLM API request failed with status: " + resp.Status + ", body: " + string(body))
	}

	// Read and process the response
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Process the response
	responseJson := string(responseBytes)

	// Parse the response
	var llmResponse LLMResponse
	if err := json.Unmarshal(responseBytes, &llmResponse); err != nil {
		return "", errors.New("Failed to parse LLM response: " + err.Error() + ", response: " + responseJson)
	}

	return llmResponse.Response, nil
}

// StreamHandler represents a function that handles streaming response chunks
type StreamHandler func(chunk string, done bool) error

// StreamGenerateContent streams the LLM responses as they are generated
func (s *LLMService) StreamGenerateContent(ctx context.Context, model string, prompt string, handler StreamHandler) error {
	// Create request with streaming enabled
	reqBody := LLMRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	// Convert request to JSON
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/generate", bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("LLM API request failed with status: " + resp.Status + ", body: " + string(body))
	}

	// Process the streaming response
	reader := bufio.NewReader(resp.Body)
	var fullResponse string

	for {
		// Read line by line (each line is a JSON object)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Skip empty lines
		if len(line) == 0 || (len(line) == 1 && line[0] == '\n') {
			continue
		}

		// Parse the JSON response
		var llmResponse LLMResponse
		if err := json.Unmarshal(line, &llmResponse); err != nil {
			return errors.New("Failed to parse LLM response chunk: " + err.Error())
		}

		// Append to full response
		fullResponse += llmResponse.Response

		// Send the chunk to the handler
		if err := handler(llmResponse.Response, llmResponse.Done); err != nil {
			return err
		}

		// If this is the last chunk, we're done
		if llmResponse.Done {
			break
		}
	}

	return nil
}
