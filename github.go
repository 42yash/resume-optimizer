package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

var tmpl = `
<form id="repo-form" class="text-sm space-y-2">
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
</form>
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

	log.Printf("Found %d repos for user %s", len(repos), username)

	t := template.Must(template.New("repos").Parse(tmpl))
	if err := t.Execute(w, repos); err != nil {
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
		return
	}
}
