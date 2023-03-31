package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Gpt4Request struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type Gpt4Response struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func callGpt4(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	url := "https://api.openai.com/v1/completions"
	requestData := Gpt4Request{
		Model:       "gpt-dv-stripe",
		Prompt:      prompt,
		MaxTokens:   512,
		Temperature: .7,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call GPT-4 API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Print the failing body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read API response body: %v", err)
		}
		fmt.Println(string(body))
		return "", fmt.Errorf("GPT-4 API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response body: %v", err)
	}

	var gpt4Response Gpt4Response
	err = json.Unmarshal(body, &gpt4Response)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal API response: %v", err)
	}

	if len(gpt4Response.Choices) > 0 {
		return gpt4Response.Choices[0].Text, nil
	}

	return "", fmt.Errorf("no choices found in GPT-4 API response")
}
