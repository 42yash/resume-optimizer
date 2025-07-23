package main

import (
	"context"
	"fmt"
	"os/exec"
)

// RepoMixOutput executes repomix with remote and stdout flags
// and returns the output as a string
func RepoMixOutput(gitRepoUrl string) (string, error) {
    // Prepare the command with repomix, --remote and --stdout flags
    cmd := exec.Command("repomix", "--remote", gitRepoUrl, "--stdout")
    fmt.Println(gitRepoUrl)
    // Capture the output
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    summary, err := summarizeUsingAI(context.Background(), string(output))
    if err != nil {
        return "", fmt.Errorf("failed to summarize content: %v", err)
    }

    // Return output as string
    return summary, nil
}

// Example usage:
// err := repomixexporter.ExportRepoAsTxt("https://github.com/username/repo")
