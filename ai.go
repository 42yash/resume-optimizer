package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

// PersonalizeResume takes a CV and job description and returns a personalized resume
func PersonalizeResume(ctx context.Context, cv, jobDescription string) (string, error) {
	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	// Create client
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create client: %v", err)
	}

	promptTemplate, err := os.ReadFile("prompt.md")
	if err != nil {
		return "", fmt.Errorf("failed to read prompt template: %v", err)
	}

	prompt := fmt.Sprintf(string(promptTemplate), cv, jobDescription)

	thinkingBudget := int32(0)

	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget,
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	fmt.Print(result.Text())

	// Extract content using the proper field access
	var fullContent strings.Builder

	if len(result.Candidates) > 0 {
		for _, part := range result.Candidates[0].Content.Parts {
			// Try direct field access first
			if part.Text != "" {
				fullContent.WriteString(part.Text)
			} else {
				// Fallback to string conversion
				fullContent.WriteString(fmt.Sprintf("%v", part))
			}
		}
	}

	return fullContent.String(), nil

}
