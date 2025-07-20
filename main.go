package main

import (
	"html/template"
	"fmt"
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ledongthuc/pdf"
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

	log.Println("Server starting on :3000...")
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
    text, err := reader.GetPlainText()
    if err != nil {
        http.Error(w, "Error extracting text from PDF", http.StatusInternalServerError)
        return
    }

    // Get job description
    jobDescription := r.FormValue("jobDescription")
    if jobDescription == "" {
        http.Error(w, "Job description is required", http.StatusBadRequest)
        return
    }

    // Print parsed content to console
    fmt.Println("=== PDF Content ===")
    fmt.Println(text)
    fmt.Println("\n=== Job Description ===")
    fmt.Println(jobDescription)

    w.Header().Set("HX-Trigger", "showMessage")
    w.Write([]byte(`
        <div id="result" class="mt-4 p-4 bg-green-100 text-green-700 rounded">
            Files received successfully! Processing...
        </div>
    `))
}
