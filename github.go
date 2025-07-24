package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func createProjectSummary(repoURLs []string) ([]string, error) {
	var wg sync.WaitGroup
	// Use buffered channels to prevent goroutines from blocking unnecessarily
	summariesChan := make(chan string, len(repoURLs))
	errorsChan := make(chan error, len(repoURLs))

	for _, repoURL := range repoURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			summary, err := RepoMixOutput(url)
			if err != nil {
				errorsChan <- err // Send the error to the errors channel
				return
			}
			summariesChan <- summary // Send the successful summary
		}(repoURL)
	}

	// Wait for all the goroutines to finish their execution
	wg.Wait()

	// Close the channels after all goroutines are done to signal we're done writing.
	close(summariesChan)
	close(errorsChan)

	// --- Robust Error Handling ---
	// Drain the errors channel and collect all errors that occurred.
	var allErrors []error
	for err := range errorsChan {
		allErrors = append(allErrors, err)
	}

	// If any errors were found, combine them into a single error message and return.
	if len(allErrors) > 0 {
		var errorMessages []string
		for _, err := range allErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return nil, fmt.Errorf("failed to process one or more repositories: %s", strings.Join(errorMessages, "; "))
	}
	
	// If there were no errors, collect all the successful summaries.
	var summaries []string
	for summary := range summariesChan {
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