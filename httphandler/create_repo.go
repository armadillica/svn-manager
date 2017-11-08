package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/armadillica/svn-manager/svnman"
)

var invalidCreatorRegexp = regexp.MustCompile(`[^\pL\d_\-., @<>'()+]+`)

// repoCreationResult is sent as JSON response to /api/repo POST requests.
type repoCreationResult struct {
	RepoID string `json:"repo_id"`
}

func (h *APIHandler) createRepo(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)

	repoInfo := svnman.CreateRepo{}
	if err := decodeJSON(w, r, &repoInfo, "create_repo", logFields); err != nil {
		return
	}

	// Convert repo ID to lower case, so limit the number of prefixes we have
	// in the repository directory.
	repoInfo.RepoID = strings.ToLower(repoInfo.RepoID)
	repoInfo.Creator = invalidCreatorRegexp.ReplaceAllString(repoInfo.Creator, " ")

	logger.Info("repository creation requested")
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	reply := repoCreationResult{RepoID: repoInfo.RepoID}
	enc := json.NewEncoder(w)
	if err := enc.Encode(reply); err != nil {
		logger.WithError(err).Error("unable to encode JSON")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to encode reply as JSON: %s", err)
		return
	}
}
