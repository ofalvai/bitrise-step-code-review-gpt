package step

import (
	"context"
	"fmt"

	"github.com/bitrise-io/go-utils/v2/log"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	apiKey string
	model  string
	logger log.Logger
}

func NewOpenAIClient(apiKey, model string, logger log.Logger) OpenAIClient {
	return OpenAIClient{
		apiKey: apiKey,
		model:  model,
		logger: logger,
	}
}

func (c OpenAIClient) GetCompletion(systemPrompt string, prPrompt string) (string, error) {
	client := openai.NewClient(c.apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prPrompt,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	c.logger.EnableDebugLog(true)
	c.logger.Printf("OpenAI API usage: %+v\n", resp.Usage)

	return resp.Choices[0].Message.Content, nil

}
