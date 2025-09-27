package http

import (
	"log"
	"net/http"

	"share-screen/pkg/infrastructure/template"
)

// StaticHandlers contains handlers for static content
type StaticHandlers struct {
	templateService *template.TemplateService
}

// NewStaticHandlers creates a new static handlers instance
func NewStaticHandlers(templateService *template.TemplateService) *StaticHandlers {
	return &StaticHandlers{
		templateService: templateService,
	}
}

// ServeIndex serves the main landing page
func (h *StaticHandlers) ServeIndex(w http.ResponseWriter, r *http.Request) {
	data := template.PageData{
		Title: "Mac → iPhone Screen Share",
	}

	if err := h.templateService.RenderPage(w, "index.html", data); err != nil {
		log.Printf("Error rendering index template: %v", err)
		http.Error(w, "Internal server error", 500)
	}
}

// ServeSender serves the sender (Mac) page
func (h *StaticHandlers) ServeSender(w http.ResponseWriter, r *http.Request) {
	data := template.PageData{
		Title:   "Sender",
		Scripts: []string{"/static/js/sender.js"},
	}

	if err := h.templateService.RenderPage(w, "sender.html", data); err != nil {
		log.Printf("Error rendering sender template: %v", err)
		http.Error(w, "Internal server error", 500)
	}
}

// ServeViewer serves the viewer (iPhone) page
func (h *StaticHandlers) ServeViewer(w http.ResponseWriter, r *http.Request) {
	data := template.PageData{
		Title:   "Viewer",
		Scripts: []string{"/static/js/viewer.js"},
	}

	if err := h.templateService.RenderPage(w, "viewer.html", data); err != nil {
		log.Printf("Error rendering viewer template: %v", err)
		http.Error(w, "Internal server error", 500)
	}
}

// ServeSenderJS serves the sender JavaScript with configured STUN server
func (h *StaticHandlers) ServeSenderJS(w http.ResponseWriter, r *http.Request) {
	data := template.PageData{}

	if err := h.templateService.RenderJS(w, "web/templates/sender.js.tmpl", data); err != nil {
		log.Printf("Error rendering sender.js template: %v", err)
		http.Error(w, "Internal server error", 500)
	}
}

// ServeViewerJS serves the viewer JavaScript with configured STUN server
func (h *StaticHandlers) ServeViewerJS(w http.ResponseWriter, r *http.Request) {
	data := template.PageData{}

	if err := h.templateService.RenderJS(w, "web/templates/viewer.js.tmpl", data); err != nil {
		log.Printf("Error rendering viewer.js template: %v", err)
		http.Error(w, "Internal server error", 500)
	}
}
