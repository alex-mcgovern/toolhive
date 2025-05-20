// Package webAssets provides access to the embedded web assets
// NOTE: go:embed can only embed assets relative to its own package, hence the
// placement in the project root directory.
package webAssets

import (
	"embed"
	"io/fs"
	"net/http"
)

// WebAssets contains the embedded web assets from the web/dist directory
//
//go:embed dist
var WebAssets embed.FS

// GetWebAssets returns an HTTP filesystem containing the embedded web assets
// It strips the "dist" prefix from the embedded files
func GetWebAssets() (http.FileSystem, error) {
	// Extract the dist subdirectory from the embedded files
	distDir, err := fs.Sub(WebAssets, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(distDir), nil
}
