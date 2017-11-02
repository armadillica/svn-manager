package svnman

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	yaml "gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

const apacheTemplate = `# Location directive for project %q
<Location /repo/%s>
    DAV svn
    SVNPath %s
    AuthType Basic
    AuthName %q
    AuthUserFile %s
    Require valid-user
</Location>
`

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
	repodir := svn.repoPath(repoInfo.RepoID)
	apafile := filepath.Join(svn.apacheConfigDir, "svn-"+repoInfo.RepoID+".conf")

	logger := log.WithFields(logFields).WithFields(log.Fields{
		"repo_id":     repoInfo.RepoID,
		"project_id":  repoInfo.ProjectID,
		"creator":     repoInfo.Creator,
		"repo_dir":    repodir,
		"apache_file": apafile,
	})

	if _, err := os.Stat(repodir); err == nil {
		logger.Warning("repository already exists")
		return ErrAlreadyExists
	}

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
		fmt.Sprintf("Blender Cloud SVN repository %q", repoInfo.RepoID),
		htpasswd)
	if err = ioutil.WriteFile(apafile, []byte(conf), 0644); err != nil {
		return err
	}

	logger.Debug("repository created, requesting Apache restart")
	svn.restarter.QueueRestart()

	return nil
}
