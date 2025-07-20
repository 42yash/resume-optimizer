go# Resume Optimizer

A web application built with Go and HTMX that helps optimize resumes based on job descriptions using AI.

## Features

- Upload resume in PDF format
- Upload job description in PDF format
- AI-powered analysis and optimization
- Simple and responsive user interface

## Prerequisites

- Go 1.21 or later
- Modern web browser

## Installation

1. Clone the repository
2. Install dependencies:
   ```
   go mod tidy
   ```

## Running the Application

1. Start the server:
   ```
   go run .
   ```
2. Open your browser and navigate to `http://localhost:3000`

## Project Structure

```
.
├── main.go              # Main application file
├── templates/           # HTML templates
│   └── index.html      # Main page template
├── static/             # Static files
└── go.mod             # Go module file
```

## Technologies Used

- Go - Backend server
- HTMX - Frontend interactivity
- TailwindCSS - Styling
- Chi Router - HTTP routing
