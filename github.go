package main

import (
	"context"
	"fmt"
	"os/exec"
)


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

func RepoMixOutput(gitRepoUrl string) (string, error) {
    // Prepare the command with repomix, --remote and --stdout flags
    cmd := exec.Command("repomix", "--remote", gitRepoUrl, "--stdout")
    fmt.Println(gitRepoUrl)
    // Capture the output
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    summary, err := summarizeRepoUsingAI(context.Background(), string(output))
    if err != nil {
        return "", fmt.Errorf("failed to summarize content: %v", err)
    }

    // Return output as string
    return summary, nil
}

