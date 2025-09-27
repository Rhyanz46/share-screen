package template

import (
	"html/template"
	"net/http"
)

// PageData represents data passed to templates
type PageData struct {
	Title      string
	ExtraHead  template.HTML
	Scripts    []string
	STUNServer string
}

// TemplateService handles template rendering
type TemplateService struct {
	templates  *template.Template
	stunServer string
}

// NewTemplateService creates a new template service
func NewTemplateService(templatesDir string, stunServer string) (*TemplateService, error) {
	tmpl, err := template.ParseGlob(templatesDir + "/*.html")
	if err != nil {
		return nil, err
	}

	return &TemplateService{
		templates:  tmpl,
		stunServer: stunServer,
	}, nil
}

// RenderPage renders a page template with base layout
func (ts *TemplateService) RenderPage(w http.ResponseWriter, templateName string, data PageData) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Set default STUN server if not provided
	if data.STUNServer == "" {
		data.STUNServer = ts.stunServer
	}

	return ts.templates.ExecuteTemplate(w, "base.html", data)
}

// RenderJS renders JavaScript template with data
func (ts *TemplateService) RenderJS(w http.ResponseWriter, templateFile string, data PageData) error {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")

	// Set default STUN server if not provided
	if data.STUNServer == "" {
		data.STUNServer = ts.stunServer
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}