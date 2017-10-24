package httphandler

import (
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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

func logFieldsForRequest(r *http.Request) (log.Fields, *log.Entry) {
	logFields := log.Fields{
		"remote_addr": r.RemoteAddr,
		"url":         r.URL,
		"method":      r.Method,
	}

	return logFields, log.WithFields(logFields)
}

// Returns the repo ID from the request, or "" when there was no (valid) one.
func getRepoID(w http.ResponseWriter, r *http.Request, logFields log.Fields) string {
	repoID, ok := mux.Vars(r)["repo-id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(logFields).Warning("no repo ID given")
		return ""
	}
	logFields["repo_id"] = repoID

	if !svnman.ValidRepoID(repoID) {
		w.WriteHeader(http.StatusBadRequest)
		log.WithFields(logFields).Warning("invalid repo ID given")
		return ""
	}
	return repoID
}

func (h *APIHandler) notImplemented(w http.ResponseWriter, r *http.Request) {
	_, logger := logFieldsForRequest(r)
	logger.Warning("handler for this URL not implemented")
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *APIHandler) createRepo(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)

	repoInfo := svnman.CreateRepo{}
	if err := decodeJSON(w, r.Body, &repoInfo, logFields); err != nil {
		return
	}

	logger.Info("going to create repository")
	if err := h.svn.CreateRepo(repoInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.WithError(err).Error("unable to create repository")
	}
}

func (h *APIHandler) modifyAccess(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)
	repoID := getRepoID(w, r, logFields)
	if repoID == "" {
		return
	}

	logger.Info("going to modify access on repository")

	h.notImplemented(w, r)
}

func (h *APIHandler) reportAccess(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}

func (h *APIHandler) deleteRepo(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}

func (h *APIHandler) blockUnblockRepo(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}

func (h *APIHandler) listAvailableHooks(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}

func (h *APIHandler) reportRepoHooks(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}

func (h *APIHandler) modifyHooks(w http.ResponseWriter, r *http.Request) {
	h.notImplemented(w, r)
}
