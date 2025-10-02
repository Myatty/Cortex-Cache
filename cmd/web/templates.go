package main

import "cortexcache.myatty.net/internal/models"

// act as holding structure for dynamic data which will be passed to html tmpls
type templateData struct {
	Snippet *models.Snippet
}
