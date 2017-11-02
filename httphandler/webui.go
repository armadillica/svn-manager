package httphandler

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// WebUI serves HTTP requests and shows a web UI.
type WebUI struct {
	root               string
	applicationVersion string
}

// TemplateData is the mapping type we use to pass data to the template engine.
type TemplateData map[string]interface{}

// CreateWebUI creates a new HTTP request handler that's bound to the given SVN Man.
func CreateWebUI(webroot, applicationVersion string) *WebUI {
	log.WithField("webroot", webroot).Info("creating web UI")
	return &WebUI{webroot, applicationVersion}
}

func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// Merges 'two' into 'one'
func merge(one map[string]interface{}, two map[string]interface{}) {
	for key := range two {
		one[key] = two[key]
	}
}

// AddRoutes adds the web UI endpoints to the router.
func (web *WebUI) AddRoutes(r *mux.Router) {
	r.HandleFunc("/", web.index).Methods("GET")

	dirname := filepath.Join(web.root, "static")
	static := noDirListing(http.StripPrefix("/static/", http.FileServer(http.Dir(dirname))))
	r.PathPrefix("/static/").Handler(static).Methods("GET")
}

func (web *WebUI) showTemplate(templfname string, w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(web.root + "/templates/" + templfname)
	if err != nil {
		_, logger := logFieldsForRequest(r)
		logger.WithError(err).WithFields(log.Fields{
			"template": templfname,
			"webroot":  web.root,
		}).Error("error parsing HTML template")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Version": web.applicationVersion,
	}

	tmpl.Execute(w, data)
}

func (web *WebUI) index(w http.ResponseWriter, r *http.Request) {
	web.showTemplate("index.html", w, r)
}
