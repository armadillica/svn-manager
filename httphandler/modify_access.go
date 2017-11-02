package httphandler

import (
	"fmt"
	"net/http"

	"github.com/armadillica/svn-manager/svnman"
)

func (h *APIHandler) modifyAccess(w http.ResponseWriter, r *http.Request) {
	logFields, logger := logFieldsForRequest(r)
	repoID := getRepoID(w, r, logFields)
	if repoID == "" {
		return
	}
	mods := svnman.ModifyAccess{}
	if err := decodeJSON(w, r, &mods, "modify_access", logFields); err != nil {
		return
	}

	logger.Info("going to modify access on repository")
	if err := h.svn.ModifyAccess(repoID, mods, logFields); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "unable to modify htpasswd: %s", err.Error())
		logger.WithError(err).Error("unable to modify htpasswd")
	}
}
