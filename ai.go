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
		You are a professional resume optimization expert. Your task is to take an existing CV and a job description, then return a highly optimized resume tailored specifically for that position.

		INSTRUCTIONS:
		- Analyze the job description to identify key requirements, skills, and keywords
		- Restructure and optimize the CV content to align with the job requirements
		- Use relevant keywords from the job description naturally throughout the resume
		- Prioritize and highlight the most relevant experiences and skills
		- Maintain truthfulness - do not fabricate information, only reorganize and emphasize existing content
		- Return ONLY the optimized resume in proper markdown format
		- Do not include any explanatory text, comments, or additional content

		INPUT FORMAT:
		**CV:**
		%s

		**Job Description:**
		%s

		OUTPUT:
		Return only the optimized resume in a clean markdown format, structured as:
		- Header with contact information
		- Professional summary/objective
		- Key skills section
		- Work experience (most relevant first)
		- Education
		- Additional relevant sections (certifications, projects, etc.)

		Focus on ATS compatibility and keyword optimization while maintaining professional formatting. Also use proper formatting such as bullet points, headings, and sections to enhance readability.`,
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

	fmt.Print("Generated resume:\n", result.Text())

	return result.Text(), nil
}
