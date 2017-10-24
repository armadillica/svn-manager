package svnman

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/armadillica/svn-manager/apache"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotImplemented is returned as error when a feature hasn't been implemented yet.
	ErrNotImplemented = errors.New("SVNMan feature not implemented")
	// ErrInvalidRepoID is returned when an invalid repository ID is used.
	ErrInvalidRepoID = errors.New("invalid repository ID given")
)

const apacheTemplate = `// Location directive for project %q
<Location /svn/%s>
    DAV svn
    SVNPath %s
    AuthType Basic
    AuthName %q
    AuthUserFile %s
    Require valid-user
</Location>
`

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

// Stored as YAML in every SVN repository we create.
type repoinfo struct {
	AppName   string    `yaml:"app_name"`
	AppVer    string    `yaml:"app_version"`
	Creation  time.Time `yaml:"created_on"`
	RepoID    string    `yaml:"repo_id"`
	ProjectID string    `yaml:"project_id"`
	Creator   string    `yaml:"creator"`
}

// CreateRepo creates a repository and Apache location directive.
func (svn *SVNMan) CreateRepo(repoInfo CreateRepo, logFields log.Fields) error {
	prefix := string([]rune(repoInfo.RepoID)[:2])
	repodir := filepath.Join(svn.repoRoot, prefix, repoInfo.RepoID)
	apafile := filepath.Join(svn.apacheConfigDir, "svn-"+repoInfo.RepoID+".conf")

	logger := log.WithFields(logFields).WithFields(log.Fields{
		"repo_id":     repoInfo.RepoID,
		"auth_realm":  repoInfo.AuthenticationRealm,
		"project_id":  repoInfo.ProjectID,
		"creator":     repoInfo.Creator,
		"repo_dir":    repodir,
		"apache_file": apafile,
	})
	logger.Info("going to create repository")

	if err := os.MkdirAll(repodir, 0750); err != nil {
		return err
	}

	// Create the SVN repository. This must happen first, because svnadmin wants the dir to be empty.
	out, err := exec.Command("svnadmin", "create", "--fs-type", "fsfs", repodir).Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			stderr := string(e.Stderr)
			logger = logger.WithField("stderr", stderr)
		}
		logger.WithError(err).Warning("error running svnadmin")
		return err
	}
	logger.WithField("stdout", string(out)).Debug("'svnadmin create' successful")

	// Create the info file.
	info := repoinfo{
		svn.appName,
		svn.appVersion,
		time.Now().UTC(),
		repoInfo.RepoID,
		repoInfo.ProjectID,
		repoInfo.Creator,
	}
	infobytes, err := yaml.Marshal(&info)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(repodir, "info.yaml"), infobytes, 0644); err != nil {
		return err
	}

	// Create an empty htpasswd file.
	htpasswd := filepath.Join(repodir, "htpasswd")
	if err = ioutil.WriteFile(htpasswd, []byte{}, 0640); err != nil {
		return err
	}

	// Create the Apache configuration file.
	conf := fmt.Sprintf(apacheTemplate,
		repoInfo.ProjectID,
		repoInfo.RepoID,
		repodir,
		repoInfo.AuthenticationRealm,
		htpasswd)
	if err = ioutil.WriteFile(apafile, []byte(conf), 0644); err != nil {
		return err
	}

	logger.Debug("repository created, requesting Apache restart")
	svn.restarter.QueueRestart()

	return nil
}

// ModifyAccess grants or revoke usage access for users on a specific repository.
func (svn *SVNMan) ModifyAccess(repoID string, mods ModifyAccess, logFields log.Fields) error {
	logger := log.WithFields(logFields).WithFields(log.Fields{
		"grant_count":  len(mods.Grant),
		"revoke_count": len(mods.Revoke),
	})

	logger.Info("modifying repository access")
	return ErrNotImplemented
}
