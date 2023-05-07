package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ozencb/tellme/cmd"
	"github.com/zalando/go-keyring"
)

const (
	APP_NAME       = "TELLMEAPP"
	OPENAI_API_URL = "https://api.openai.com/v1/chat/completions"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {
	cmd.Execute()

	apiKey, err := getOrSetApiKey()
	if err != nil {
		fmt.Println(err)
	}

	prompt := os.Args[1]

	request := OpenAIRequest{
		Model: cmd.Model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	response, err := makeRequest(request, apiKey)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Message.Content)
	}
}

func getOrSetApiKey() (string, error) {
	apiKey, err := keyring.Get(APP_NAME, "API_KEY")

	if cmd.ResetApiKey || err != nil {
		fmt.Println("Enter your OpenAI API Key: ")
		fmt.Scanln(&apiKey)
		err := keyring.Set(APP_NAME, "API_KEY", apiKey)
		if err != nil {
			fmt.Println(err)
		}
		return apiKey, nil
	}
	return apiKey, nil
}

func makeRequest(request OpenAIRequest, apiKey string) (OpenAIResponse, error) {
	payload, err := json.Marshal(request)

	if err != nil {
		return OpenAIResponse{}, err
	}

	req, err := http.NewRequest("POST", OPENAI_API_URL, bytes.NewBuffer(payload))
	if err != nil {
		return OpenAIResponse{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return OpenAIResponse{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OpenAIResponse{}, err
	}

	var response OpenAIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return OpenAIResponse{}, err
	}

	return response, nil
}
