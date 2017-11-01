package svnman

import (
	"path/filepath"

	"github.com/foomo/htpasswd"
	log "github.com/sirupsen/logrus"
)

// ModifyAccess grants or revoke usage access for users on a specific repository.
func (svn *SVNMan) ModifyAccess(repoID string, mods ModifyAccess, logFields log.Fields) error {
	filename := filepath.Join(svn.repoPath(repoID), "htpasswd")
	logger := log.WithFields(logFields).WithFields(log.Fields{
		"grant_count":  len(mods.Grant),
		"revoke_count": len(mods.Revoke),
		"filename":     filename,
	})

	logger.Debug("modifying repository access")
	passwds, err := htpasswd.ParseHtpasswdFile(filename)
	if err != nil {
		logger.WithError(err).Error("unable to parse htpasswd")
		return err
	}

	for _, grant := range mods.Grant {
		passwds[grant.Username] = grant.Password
	}
	for _, revoke := range mods.Revoke {
		delete(passwds, revoke)
	}

	if err := passwds.WriteToFile(filename); err != nil {
		logger.WithError(err).Error("unable to save htpasswd")
		return err
	}

	logger.Info("repository access modified")
	return nil
}
