package svnman

import (
	"errors"
	"path/filepath"

	"github.com/armadillica/svn-manager/apache"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotImplemented is returned as error when a feature hasn't been implemented yet.
	ErrNotImplemented = errors.New("SVNMan feature not implemented")
	// ErrInvalidRepoID is returned when an invalid repository ID is used.
	ErrInvalidRepoID = errors.New("invalid repository ID given")
	// ErrAlreadyExists is returned when a request to create a repository fails because it already exists.
	ErrAlreadyExists = errors.New("repository with this ID already exists")
)

// Manager contains the interface of SVNMan, for testing/mocking purposes.
type Manager interface {
	CreateRepo(repoInfo CreateRepo, logFields log.Fields) error
	ModifyAccess(repoID string, mods ModifyAccess, logFields log.Fields) error
}

// SVNMan provides SVN management operations.
type SVNMan struct {
	restarter       apache.Restarter
	repoRoot        string
	apacheConfigDir string

	// To store in the info.txt file.
	appName    string
	appVersion string
}

// Create returns a newly created SVNMan instance.
func Create(restarter apache.Restarter, repoRoot, apacheConfigDir, appName, appVersion string) *SVNMan {
	log.WithFields(log.Fields{
		"repo_root": repoRoot,
		"apache":    apacheConfigDir,
	}).Info("creating SVN manager")
	return &SVNMan{restarter, repoRoot, apacheConfigDir, appName, appVersion}
}

func (svn *SVNMan) repoPath(repoID string) string {
	prefix := string([]rune(repoID)[:2])
	return filepath.Join(svn.repoRoot, prefix, repoID)
}
