package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

var tmpl = `
  {{range .}}
    <div>
      <label class="inline-flex items-center">
        <input type="checkbox" name="repos" value="{{.Name}}" class="mr-2">
        {{.Name}}
      </label>
    </div>
  {{else}}
    <p>No repositories found.</p>
  {{end}}
`

func handleRepos(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("githubUsername")
	log.Printf("Fetching repos for user: %s", username)
	if username == "" {
		http.Error(w, "Username required", http.StatusBadRequest)
		return
	}

	apiURL := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Could not fetch GitHub repos", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "GitHub API error", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	var repos []Repo
	if err := json.Unmarshal(body, &repos); err != nil {
		http.Error(w, "Failed to decode repos", http.StatusInternalServerError)
		return
	}

	// Convert repo names to URLs
	for i := range repos {
		repos[i].URL = fmt.Sprintf("https://github.com" + username + "/" + repos[i].Name)
	}

	log.Printf("Found %d repos for user %s", len(repos), username)

	t := template.Must(template.New("repos").Parse(tmpl))
	if err := t.Execute(w, repos); err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}
}

func createProjectSummary(repoURLs []string) ([]string, error) {
    var summaries []string

    for _, repoURL := range repoURLs {
        summary, err := RepoMixOutput(repoURL)
        if err != nil {
            return nil, err
        }

        summaries = append(summaries, summary)
    }
    return summaries, nil
}

func summarizeUsingAI(ctx context.Context, content string) (string, error) {
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

	promptTemplate, err := os.ReadFile("repo-prompt.md")
	if err != nil {
		return "", fmt.Errorf("failed to read prompt template: %v", err)
	}

	prompt := fmt.Sprintf(string(promptTemplate), content)

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

	// fmt.Println("Generated summary:", fullContent.String())

	return fullContent.String(), nil
}