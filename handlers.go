package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/ledongthuc/pdf"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "index.html")))
	tmpl.Execute(w, nil)
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form with 10 MB max memory
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
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

	ctx := context.Background()

	// Generate personalized resume
	optimizedResume, err := PersonalizeResume(ctx, originalResume, jobDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	data := struct {
		OriginalResume  string
		JobDescription  string
		OptimizedResume string
	}{
		OriginalResume:  template.HTMLEscapeString(originalResume),
		JobDescription:  template.HTMLEscapeString(jobDescription),
		OptimizedResume: template.JSEscapeString(optimizedResume),
	}

	// Render the result template
	tmpl := template.Must(template.ParseFiles(filepath.Join("templates", "result.html")))
	err = tmpl.Execute(w, data)
	if err != nil {
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
