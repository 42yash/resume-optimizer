# Resume Optimizer

A web application built with Go and HTMX that helps optimize resumes based on job descriptions using Google's Gemini AI model.

## Features

- Upload resume in PDF format
- Enter job description text
- AI-powered resume analysis and optimization
- Markdown output with copy functionality
- Simple and responsive user interface using TailwindCSS
- ATS-friendly resume formatting
- Real-time processing feedback

## Prerequisites

- Go 1.21 or later
- Modern web browser
- Google Gemini API key

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up environment variables:
   ```bash
   export GEMINI_API_KEY=your_api_key_here
   ```

## Running the Application

1. Start the server:
   ```bash
   go run .
   ```
2. Open your browser and navigate to `http://localhost:3000`

## Project Structure

```
.
├── main.go              # Server setup and routing configuration
├── handlers.go          # HTTP request handlers
├── ai.go               # Gemini AI integration
├── templates/          # HTML templates
│   ├── index.html     # Upload form template
│   └── result.html    # Results page template
├── static/            # Static assets directory
├── go.mod            # Go module definition
└── go.sum            # Go module checksums
```

## Technologies Used

- Go - Backend server and PDF processing
- Google Gemini AI - Resume optimization
- HTMX - Frontend interactivity
- TailwindCSS - Responsive styling
- Chi Router - HTTP routing and middleware
- marked.js - Markdown rendering
- ledongthuc/pdf - PDF text extraction

## API Endpoints

- `GET /` - Serves the main upload form
- `POST /process` - Handles resume processing
  - Accepts multipart form data with:
    - `resume` (PDF file)
    - `jobDescription` (text)

## How It Works

1. User uploads their resume in PDF format
2. User provides the target job description
3. Backend extracts text from the PDF
4. Gemini AI analyzes and optimizes the resume for the specific job
5. Results are displayed in a formatted markdown view with copy functionality

## Dependencies

- github.com/go-chi/chi/v5 - HTTP routing
- github.com/go-chi/cors - CORS middleware
- google.golang.org/genai - Gemini AI client
- github.com/ledongthuc/pdf - PDF processing

## Error Handling

- Validates PDF file uploads
- Checks for required environment variables
- Provides user-friendly error messages
- Handles API rate limiting and timeouts

## Security Features

- CORS configuration
- File size limits (10MB)
- Content type validation
- Secure template rendering

## Browser Compatibility

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

This project is open source and available