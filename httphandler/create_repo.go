package httphandler

import (
	"fmt"
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
)

func (h *APIHandler) createRepo(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)

	repoInfo := svnman.CreateRepo{}
	if err := decodeJSON(w, r, &repoInfo, "create_repo", logFields); err != nil {
		return
	}

	logger.Info("going to create repository")
	err := h.svn.CreateRepo(repoInfo, logFields)
	if err == svnman.ErrAlreadyExists {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "repository %q already exists", repoInfo.RepoID)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to create repository: %s", err.Error())
		logger.WithError(err).Error("unable to create repository")
		return
	}

	route, err := h.r.Get("get-repo").URL("repo-id", repoInfo.RepoID)
	if err != nil {
		logger.WithError(err).WithField("repo_id", repoInfo.RepoID).Error("unable to find URL for repository")
	} else {
		w.Header().Set("Location", route.String())
	}
	w.WriteHeader(http.StatusCreated)
}
