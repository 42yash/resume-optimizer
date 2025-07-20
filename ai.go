package main

import (
	"context"
	"fmt"
	"os"

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

	// Create the prompt
	prompt := fmt.Sprintf(`
		Please personalize this resume for the specific job description provided. 
		Tailor the content to highlight relevant skills, experience, and achievements that match the job requirements.
		Keep the same format but emphasize the most relevant aspects.

		Original CV:
		%s

		Job Description:
		%s

		Please provide the personalized resume:`,
		cv,
		jobDescription,
	)

	thinkingBudget := int32(0)

	// Call Gemini API
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget: &thinkingBudget,
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	return result.Text(), nil
}
