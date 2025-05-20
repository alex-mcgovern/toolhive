package app

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	s "github.com/stacklok/toolhive/pkg/api"
	webAssets "github.com/stacklok/toolhive/web"
)

// Port for serving the static web interface
var webPort int

// Constants for web command
const (
	defaultHost    = "127.0.0.1"
	defaultPort    = 8080
	defaultWebPort = 8081
)

func init() {
	webCmd.Flags().StringVar(&host, "host", defaultHost, "Host address to bind the servers to")
	webCmd.Flags().IntVar(&port, "port", defaultPort, "Port to bind the API server to")
	webCmd.Flags().IntVar(&webPort, "web-port", defaultWebPort, "Port to bind the web interface server to")
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the ToolHive web interface",
	Long:  `Starts the ToolHive API server with a web interface and opens your default browser.`,
	RunE:  webCmdFunc,
}

// webCmdFunc implements the main functionality for the web command
func webCmdFunc(cmd *cobra.Command, _ []string) error {
	// Ensure server is shutdown gracefully on Ctrl+C
	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
	defer cancel()

	// Get debug mode flag
	debugMode, _ := cmd.Flags().GetBool("debug")

	// Format addresses for the API and web servers
	apiAddress := fmt.Sprintf("%s:%d", host, port)
	webAddress := fmt.Sprintf("%s:%d", host, webPort)

	// Find the web directory containing the Vite React application
	webDir, err := findWebDirectory()
	if err != nil {
		return err
	}

	// Start both servers and wait for completion
	return startServersAndWait(ctx, apiAddress, webAddress, webDir, debugMode)
}

// findWebDirectory attempts to locate the web/dist directory containing the Vite React application,
// either as embedded assets or in the file system.
func findWebDirectory() (string, error) {
	// First try to use embedded assets
	if embeddedFS, err := useEmbeddedAssets(); err == nil {
		return embeddedFS, nil
	}

	// Fallback to file system for development mode
	// In production, we expect the web/dist directory to be relative to the executable
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	// Check the production/installed location
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(execPath)))
	prodWebDir := filepath.Join(rootDir, "web", "dist")
	if _, err := os.Stat(prodWebDir); err == nil {
		return prodWebDir, nil
	}

	// For development environments, check the current directory
	cwd, err := os.Getwd()
	if err == nil {
		devWebDir := filepath.Join(cwd, "web", "dist")
		if _, err := os.Stat(devWebDir); err == nil {
			return devWebDir, nil
		}
	}

	// Provide a clear error message if the directory wasn't found
	return "", fmt.Errorf("could not find web/dist directory - make sure you've built the frontend with 'npm run build' in the web directory")
}

// startServersAndWait starts both the API and web servers in goroutines,
// opens the browser pointing to the web interface, and waits for completion
// or interruption signals
func startServersAndWait(ctx context.Context, apiAddress, webAddress, webDir string, debugMode bool) error {
	errChanAPI := startAPIServer(ctx, apiAddress, debugMode)
	errChanWeb := startWebServer(webAddress, webDir)

	openBrowserToWebInterface(webAddress)

	// Wait for either server to return an error or for interrupt signal
	return waitForCompletion(ctx, errChanAPI, errChanWeb)
}

// startAPIServer starts the API server in a goroutine and returns a channel for errors
func startAPIServer(ctx context.Context, apiAddress string, debugMode bool) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting API server on %s\n", apiAddress)
		errChan <- s.Serve(ctx, apiAddress, debugMode)
	}()
	return errChan
}

// startWebServer starts the web server in a goroutine and returns a channel for errors
func startWebServer(webAddress, webDir string) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("Starting web server on %s and serving content from %s\n", webAddress, webDir)

		// Create the web server with SPA routing support
		server := &http.Server{
			Addr:              webAddress,
			Handler:           createSPAHandler(webDir),
			ReadHeaderTimeout: 10 * time.Second,
		}

		errChan <- server.ListenAndServe()
	}()
	return errChan
}

// createSPAHandler creates an HTTP handler that supports Single Page
// Application routing.

// This implements the "history API fallback" pattern required for React Router
// to work properly:
//   - Static files (JS, CSS, images) are served directly when they exist
//   - All other requests that would 404 (like /foo) serve index.html instead,
//     allowing React Router to handle client-side routing based on the URL path
//
// Without this handler, direct navigation to routes like /foo would fail
// because the server would look for a /foo file that doesn't exist.
func createSPAHandler(webDir string) http.Handler {
	// Check if we're using embedded assets
	if webDir == "::embedded::" {
		return createEmbeddedSPAHandler()
	}

	// Create a file server handler for static assets on disk
	fileServer := http.FileServer(http.Dir(webDir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested path exists as a file
		path := filepath.Join(webDir, r.URL.Path)
		_, err := os.Stat(path)

		// Serve existing files directly
		if err == nil || !os.IsNotExist(err) {
			fileServer.ServeHTTP(w, r)
			return
		}

		// For all other paths, serve index.html for React Router
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})
}

// createEmbeddedSPAHandler creates an HTTP handler for embedded assets with SPA support
func createEmbeddedSPAHandler() http.Handler {
	// Get the embedded filesystem
	webFS, err := webAssets.GetWebAssets()
	if err != nil {
		// If there's an error, return a handler that always returns an error
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Embedded assets not available: "+err.Error(), http.StatusInternalServerError)
		})
	}

	// Create a file server for the embedded assets
	fileServer := http.FileServer(webFS)

	// Return a handler that implements the SPA pattern
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly first
		origPath := r.URL.Path

		// Open the file to check if it exists
		f, err := webFS.Open(origPath)
		if err == nil {
			f.Close() // File exists, close it and serve normally
			fileServer.ServeHTTP(w, r)
			return
		}

		// If file doesn't exist, serve index.html for the SPA
		// Reset the URL path to serve the index.html file
		r2 := new(http.Request)
		*r2 = *r // Clone the request
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = "index.html"
		fileServer.ServeHTTP(w, r2)
	})
}

// useEmbeddedAssets returns a special path indicating we're using embedded assets.
// This is used as a signal to createSPAHandler to use the embedded filesystem.
func useEmbeddedAssets() (string, error) {
	// This special string indicates we're using embedded assets
	// It doesn't represent a real file path and is used internally
	return "::embedded::", nil
}

// openBrowser opens the specified URL in the default browser of the user's OS
func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

// openBrowserToWebInterface opens the user's default browser to the web interface
func openBrowserToWebInterface(webAddress string) {
	url := fmt.Sprintf("http://%s", webAddress)
	fmt.Printf("Opening %s in your default browser...\n", url)
	if err := openBrowser(url); err != nil {
		fmt.Printf("Error opening browser: %v\n", err)
		fmt.Printf("Please open %s manually in your browser\n", url)
	}
}

// waitForCompletion waits for either server to return an error or for an interrupt signal
func waitForCompletion(ctx context.Context, errChanAPI, errChanWeb <-chan error) error {
	select {
	case err := <-errChanAPI:
		return fmt.Errorf("API server error: %w", err)
	case err := <-errChanWeb:
		return fmt.Errorf("Web server error: %w", err)
	case <-ctx.Done():
		fmt.Println("Servers shutting down...")
		return nil
	}
}
