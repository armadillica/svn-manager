package httphandler

import (
	"fmt"
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
)

func (h *APIHandler) deleteRepo(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)
	repoID := getRepoID(w, r, logFields)
	if repoID == "" {
		return
	}
	logger = logger.WithField("repo_id", repoID)
	logger.Info("repository deletion requested")

	err := h.svn.DeleteRepo(repoID, logFields)
	if err == svnman.ErrNotFound {
		logger.Warning("nonexistent repository requested")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "nonexistent repository requested")
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to delete repository")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to delete repository: %s", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
