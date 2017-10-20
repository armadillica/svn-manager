package httphandler

import (
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
	"github.com/gorilla/mux"
)

// APIHandler serves HTTP requests and forwards connections to the SVN Man.
type APIHandler struct {
	svn *svnman.SVNMan
}

// CreateHTTPHandler creates a new HTTP request handler that's bound to the given SVN Man.
func CreateHTTPHandler(svn *svnman.SVNMan) *APIHandler {
	return &APIHandler{svn}
}

// AddRoutes adds the web endpoints to the router.
func (h *APIHandler) AddRoutes(r *mux.Router) {
	r.HandleFunc("/repo", h.createRepo).Methods("POST")
	r.HandleFunc("/repo/{repo-id}", h.deleteRepo).Methods("DELETE")
	r.HandleFunc("/repo/{repo-id}/block", h.blockUnblockRepo).Methods("POST")
	r.HandleFunc("/repo/{repo-id}/access", h.modifyAccess).Methods("POST")
	r.HandleFunc("/repo/{repo-id}/access", h.reportAccess).Methods("GET")
	r.HandleFunc("/repo/{repo-id}/hooks", h.reportRepoHooks).Methods("GET")
	r.HandleFunc("/repo/{repo-id}/hooks", h.modifyHooks).Methods("POST")
	r.HandleFunc("/hooks", h.listAvailableHooks).Methods("GET")
}

func (h *APIHandler) createRepo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) modifyAccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) reportAccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) deleteRepo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) blockUnblockRepo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) listAvailableHooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) reportRepoHooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) modifyHooks(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
