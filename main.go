package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ledongthuc/pdf"

	"google.golang.org/genai"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static files
	fileServer := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Routes
	r.Get("/", handleHome)
	r.Post("/process", handleProcess)

	log.Println("Server starting on http://localhost:3000...")
	http.ListenAndServe(":3000", r)
}

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

	// Get the PDF file (removed fileHeader variable)
	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving resume file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the PDF content
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		http.Error(w, "Error reading PDF file", http.StatusInternalServerError)
		return
	}

	// Parse PDF content
	reader, err := pdf.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		http.Error(w, "Error parsing PDF", http.StatusInternalServerError)
		return
	}

	// Extract text from PDF
	cvReader, err := reader.GetPlainText()
	if err != nil {
		http.Error(w, "Error extracting text from PDF", http.StatusInternalServerError)
		return
	}

	// Convert io.Reader to string
	cvBytes, err := io.ReadAll(cvReader)
	if err != nil {
		http.Error(w, "Error reading PDF text content", http.StatusInternalServerError)
		return
	}
	cvtext := string(cvBytes)

	// Get job description
	jdtext := r.FormValue("jobDescription")
	if jdtext == "" {
		http.Error(w, "Job description is required", http.StatusBadRequest)
		return
	}

	// Print parsed content to console
	fmt.Println("=== PDF Content ===")
	fmt.Println(cvtext)
	fmt.Println("\n=== Job Description ===")
	fmt.Println(jdtext)

	// TODO: Add AI processing logic here
	// Sample placeholder for processing logic
	ctx := context.Background()

	// Generate personalized resume

	personalizedResume, err := PersonalizeResume(ctx, cvtext, jdtext)
	if err != nil {
		log.Fatal(err)
	}

	// For now, just return a success message
	w.Header().Set("HX-Trigger", "showMessage")
	w.Write([]byte(`
		<div id="result" class="mt-4 p-4 bg-green-100 text-green-700 rounded">
			<h2 class="text-lg font-semibold">Personalized Resume Generated</h2>

			<p class="mt-2">Your personalized resume is ready:</p>

			<pre class="mt-2 bg-white p-4 rounded shadow">
			` + personalizedResume + `
			</pre>
		</div>
	`))
}

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
