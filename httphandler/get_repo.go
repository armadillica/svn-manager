package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
)

// RepoDescription is sent as JSON response to /api/repo/{repo-id} requests.
type RepoDescription struct {
	RepoID string   `json:"repo_id"`
	Access []string `json:"access"` // list of usernames
}

func (h *APIHandler) getRepo(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)
	repoID := getRepoID(w, r, logFields)
	if repoID == "" {
		return
	}
	logger = logger.WithField("repo_id", repoID)

	names, err := h.svn.GetUsernames(repoID)
	if err == svnman.ErrNotFound {
		logger.Warning("nonexistent repository requested")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "nonexistent repository requested")
		return
	} else if err != nil {
		logger.WithError(err).Error("unable to get usernames for repo")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to get usernames: %s", err)
		return
	}

	reply := RepoDescription{
		RepoID: repoID,
		Access: names,
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(reply); err != nil {
		logger.WithError(err).Error("unable to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode reply as JSON: %s", err)
		return
	}

	return
}
