package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	tmpl.Execute(w, nil)
}

func handleLinkedInVerify(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("linkedinUrl")
	if url == "" || !strings.Contains(url, "linkedin.com/in/") {
		fmt.Fprintf(w, "LinkedIn profile is invalid!")
	}
	fmt.Fprintf(w, "✅ LinkedIn profile appears valid!")
}

func handleProcess(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	r = r.WithContext(ctx)

	// Parse the multipart form with 10 MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the PDF file
	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving resume file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Convert PDF to text
	originalResume, err := convertPDFToText(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get job description
	jobDescription := r.FormValue("jobDescription")
	if jobDescription == "" {
		http.Error(w, "Job description is required", http.StatusBadRequest)
		return
	}

	selectedRepos := r.Form["repos"]
	//  Convert repo name to url
	if len(selectedRepos) == 0 {
		http.Error(w, "At least one repository must be selected", http.StatusBadRequest)
		return
	}	

	username := r.FormValue("githubUsername")

	for i, repoName := range selectedRepos {
		repoURL := fmt.Sprintf("https://github.com/" + username + "/" + repoName)
		selectedRepos[i] = repoURL
	}

	projectSummaries, err := createProjectSummary(selectedRepos)
	if err != nil {
		http.Error(w, "Error creating project summaries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Selected repositories:", selectedRepos)

	// Generate personalized resume in markdown format
	optimizedResume, err := PersonalizeResume(r.Context(), originalResume, jobDescription, projectSummaries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the result template with just the optimized resume
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "resume.html")))
	if err := tmpl.Execute(w, struct{ OptimizedResume string }{optimizedResume})   ; err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func convertPDFToText(file io.Reader) (string, error) {
	// Read the PDF content
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", fmt.Errorf("error reading PDF file: %v", err)
	}

	// Parse PDF content
	reader, err := pdf.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return "", fmt.Errorf("error parsing PDF: %v", err)
	}

	// Extract text from PDF
	cvReader, err := reader.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("error extracting text from PDF: %v", err)
	}

	// Convert io.Reader to string
	cvBytes, err := io.ReadAll(cvReader)
	if err != nil {
		return "", fmt.Errorf("error reading PDF text content: %v", err)
	}

	return string(cvBytes), nil
}
