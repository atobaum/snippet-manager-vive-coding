package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/atobaum/snippet-manager/internal/snippet"
)

//go:embed all:dist
var staticFiles embed.FS

// Server represents the web server
type Server struct {
	snippetService *snippet.Service
	port           int
	devMode        bool
}

// NewServer creates a new web server
func NewServer(port int, devMode bool) (*Server, error) {
	svc, err := snippet.NewService()
	if err != nil {
		return nil, fmt.Errorf("failed to create snippet service: %w", err)
	}

	return &Server{
		snippetService: svc,
		port:           port,
		devMode:        devMode,
	}, nil
}

// Start starts the web server
func (s *Server) Start() error {
	// Set up routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/snippets", s.handleSnippets)
	mux.HandleFunc("/api/snippets/", s.handleSnippet)

	// Static files handling
	if s.devMode {
		// In dev mode, proxy to Svelte dev server
		s.setupDevProxy(mux)
	} else {
		// In production mode, serve embedded files
		s.setupStaticFiles(mux)
	}

	fmt.Printf("🚀 Server starting on http://localhost:%d\n", s.port)
	if s.devMode {
		fmt.Println("📝 Development mode: Make sure Svelte dev server is running on port 5173")
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.corsMiddleware(mux))
}

// setupDevProxy sets up proxy to Svelte dev server for development
func (s *Server) setupDevProxy(mux *http.ServeMux) {
	target, err := url.Parse("http://localhost:5173")
	if err != nil {
		fmt.Printf("Error parsing dev server URL: %v\n", err)
		s.setupStaticFiles(mux)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		proxy.ServeHTTP(w, r)
	})
}

// setupStaticFiles sets up static file serving for production
func (s *Server) setupStaticFiles(mux *http.ServeMux) {
	staticFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		// If dist folder doesn't exist, serve a simple message
		mux.HandleFunc("/", s.handleFallback)
		return
	}

	// Serve static files with SPA routing support
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to open the requested file
		file, err := staticFS.Open(path)
		if err != nil {
			// If file not found and not a static asset, serve index.html for SPA routing
			if os.IsNotExist(err) && !strings.Contains(path, ".") {
				indexFile, err := staticFS.Open("index.html")
				if err != nil {
					http.Error(w, "index.html not found", http.StatusNotFound)
					return
				}
				defer indexFile.Close()

				w.Header().Set("Content-Type", "text/html")
				io.Copy(w, indexFile)
				return
			}
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		// Set appropriate content type
		if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Content-Type", "text/html")
		}

		io.Copy(w, file)
	})
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleSnippets handles GET /api/snippets and POST /api/snippets
func (s *Server) handleSnippets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSnippets(w, r)
	case http.MethodPost:
		s.createSnippet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleSnippet handles GET/PUT/DELETE /api/snippets/{name}
func (s *Server) handleSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract snippet name from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/snippets/")
	name := strings.Split(path, "/")[0]

	if name == "" {
		http.Error(w, "Snippet name required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getSnippet(w, r, name)
	case http.MethodPut:
		s.updateSnippet(w, r, name)
	case http.MethodDelete:
		s.deleteSnippet(w, r, name)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getSnippets returns all snippets
func (s *Server) getSnippets(w http.ResponseWriter, r *http.Request) {
	snippets, err := s.snippetService.ListSnippets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippets)
}

// createSnippet creates a new snippet
func (s *Server) createSnippet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Command     string   `json:"command"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	err := s.snippetService.CreateSnippet(req.Name, req.Description, req.Command, req.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Snippet created successfully"})
}

// getSnippet returns a specific snippet
func (s *Server) getSnippet(w http.ResponseWriter, r *http.Request, name string) {
	snippet, err := s.snippetService.GetSnippet(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}

// updateSnippet updates an existing snippet
func (s *Server) updateSnippet(w http.ResponseWriter, r *http.Request, name string) {
	var req struct {
		Description string   `json:"description"`
		Command     string   `json:"command"`
		Tags        []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := s.snippetService.UpdateSnippet(name, req.Description, req.Command, req.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Snippet updated successfully"})
}

// deleteSnippet deletes a snippet
func (s *Server) deleteSnippet(w http.ResponseWriter, r *http.Request, name string) {
	err := s.snippetService.DeleteSnippet(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Snippet deleted successfully"})
}

// handleFallback serves a fallback page when dist folder doesn't exist
func (s *Server) handleFallback(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Snippet Manager</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 600px; margin: 0 auto; }
        .status { background: #f0f0f0; padding: 20px; border-radius: 8px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚧 Snippet Manager Web UI</h1>
        <div class="status">
            <p><strong>Status:</strong> Server is running!</p>
            <p><strong>API Endpoints:</strong></p>
            <ul>
                <li>GET /api/snippets - List all snippets</li>
                <li>POST /api/snippets - Create new snippet</li>
                <li>GET /api/snippets/{name} - Get specific snippet</li>
                <li>PUT /api/snippets/{name} - Update snippet</li>
                <li>DELETE /api/snippets/{name} - Delete snippet</li>
            </ul>
            <p><em>Web UI is coming soon... Build the Svelte app first!</em></p>
        </div>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
