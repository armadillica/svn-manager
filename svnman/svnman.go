package svnman

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/armadillica/svn-manager/apache"
	"github.com/foomo/htpasswd"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotImplemented is returned as error when a feature hasn't been implemented yet.
	ErrNotImplemented = errors.New("SVNMan feature not implemented")
	// ErrInvalidRepoID is returned when an invalid repository ID is used.
	ErrInvalidRepoID = errors.New("invalid repository ID given")
	// ErrAlreadyExists is returned when a request to create a repository fails because it already exists.
	ErrAlreadyExists = errors.New("repository with this ID already exists")
	// ErrNotFound indicates that the requested repository does not exist.
	ErrNotFound = errors.New("repository with this ID does not exist")
	// ErrDeletion indicates that a repository deletion failed. Specifics are logged.
	ErrDeletion = errors.New("unable to delete repository")
)

// RFC3339fs is a filesystem-friendly version of RFC3339.
const RFC3339fs = "2006-01-02T15-04-05Z07-00"

// Manager contains the interface of SVNMan, for testing/mocking purposes.
type Manager interface {
	CreateRepo(repoInfo CreateRepo, logFields log.Fields) error
	ModifyAccess(repoID string, mods ModifyAccess, logFields log.Fields) error
	GetUsernames(repoID string) ([]string, error)
	DeleteRepo(repoID string, logFields log.Fields) error
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

func (svn *SVNMan) atticPath(repoID string, timestamp time.Time) string {
	prefix := string([]rune(repoID)[:2])
	fname := repoID + "-" + timestamp.Format(RFC3339fs)
	return filepath.Join(svn.repoRoot, "attic", prefix, fname)
}

func (svn *SVNMan) apaConfPath(repoID string) string {
	prefix := string([]rune(repoID)[:2])
	fname := "svn-" + repoID + ".conf"
	return filepath.Join(svn.apacheConfigDir, prefix, fname)
}

func (svn *SVNMan) apaAtticPath(repoID string, timestamp time.Time) string {
	prefix := string([]rune(repoID)[:2])
	fname := "svn-" + repoID + ".conf-" + timestamp.Format(RFC3339fs)
	return filepath.Join(svn.apacheConfigDir, "attic", prefix, fname)
}

func (svn *SVNMan) htpasswd(repoID string) string {
	return filepath.Join(svn.repoPath(repoID), "htpasswd")
}

// GetUsernames returns the list of usernames that have access to the given repository.
func (svn *SVNMan) GetUsernames(repoID string) ([]string, error) {
	filename := svn.htpasswd(repoID)
	passwds, err := htpasswd.ParseHtpasswdFile(filename)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	i := 0
	names := make([]string, len(passwds))
	for name := range passwds {
		names[i] = name
		i++
	}

	return names, nil
}
