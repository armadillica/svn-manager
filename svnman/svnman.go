package svnman

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// ErrNotImplemented is returned as error when a feature hasn't been implemented yet.
var ErrNotImplemented = errors.New("SVNMan feature not implemented")

// SVNMan provides SVN management operations.
type SVNMan struct {
}

// CreateRepo creates a repository and Apache location directive.
func (svn *SVNMan) CreateRepo(repoInfo CreateRepo, logFields log.Fields) error {
	logger := log.WithFields(logFields).WithFields(log.Fields{
		"repo_id":       repoInfo.RepoID,
		"auth_realm":    repoInfo.AuthenticationRealm,
		"project_id":    repoInfo.ProjectID,
		"creator_name":  repoInfo.Creator.FullName,
		"creator_email": repoInfo.Creator.Email,
	})

	logger.Info("going to create repository")
	return ErrNotImplemented
}
}
