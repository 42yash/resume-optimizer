package repomixexporter

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ExportRepoAsTxt fetches a GitHub repo, extracts code with repomix, writes to txt, and cleans up.
func ExportRepoAsTxt(repoURL string) error {
	// Fix output path
	outputPath := "output.txt"

    // Construct the repomix command:
    // repomix --remote <repoURL> --style plain
    cmd := exec.Command("repomix", "--remote", repoURL, "--style", "plain")

    outFile, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer outFile.Close()

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }

    stderr, err := cmd.StderrPipe()
    if err != nil {
        return err
    }

    if err := cmd.Start(); err != nil {
        return err
    }

    // Write cleaned output from stdout to file
    cleanWrite(stdout, outFile)

    // Optionally, handle stderr or log warnings
    errs, _ := io.ReadAll(stderr)
    if len(errs) > 0 {
        return errors.New(string(errs))
    }

    // Wait for finish
    return cmd.Wait()
}

// cleanWrite removes excess whitespace from output and writes to dst.
func cleanWrite(src io.Reader, dst io.Writer) {
    scanner := bufio.NewScanner(src)
    writer := bufio.NewWriter(dst)
    defer writer.Flush()
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if len(line) > 0 {
            writer.WriteString(line + "\n")
        }
    }
}

// Example usage:
// err := repomixexporter.ExportRepoAsTxt("https://github.com/username/repo")
