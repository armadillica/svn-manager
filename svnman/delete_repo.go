package svnman

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// DeleteRepo moves a repository into the attic.
func (svn *SVNMan) DeleteRepo(repoID string, logFields log.Fields) error {
	logger := log.WithFields(logFields)
	logger.Debug("deleting repository")

	timestamp := time.Now()
	apaConfPath := svn.apaConfPath(repoID)
	apaAtticPath := svn.apaAtticPath(repoID, timestamp)
	repoPath := svn.repoPath(repoID)
	atticPath := svn.atticPath(repoID, timestamp)

	logger = logger.WithFields(log.Fields{
		"apa_conf":  apaConfPath,
		"apa_attic": apaAtticPath,
		"repo":      repoPath,
		"attic":     atticPath,
	})

	if err := os.MkdirAll(filepath.Dir(atticPath), 0750); err != nil {
		logger.WithError(err).Error("unable to create attic path for repo")
		return ErrDeletion
	}
	if err := os.MkdirAll(filepath.Dir(apaAtticPath), 0750); err != nil {
		logger.WithError(err).Error("unable to create attic path for Apache config")
		return ErrDeletion
	}

	// Remove the Apache configuration first. With that, the repository
	// should be inaccessible, even when the repo files themselves cannot
	// be moved to the attic. Doing it the other way around (repo dir first)
	// will cause errors when the Apache file cannot be moved but the repo can.
	err := os.Rename(apaConfPath, apaAtticPath)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.WithError(err).Error("unable to move Apache config to attic")
			return ErrDeletion
		}
		logger.Warning("trying to remove non-existant Apache config file")
	}

	err = os.Rename(repoPath, atticPath)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.WithError(err).Error("unable to move repository to attic")
			return ErrDeletion
		}
		logger.Warning("trying to remove non-existant repository")
	}

	svn.restarter.QueueRestart()
	logger.Info("repository deleted")
	return nil
}
