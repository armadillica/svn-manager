package svnman

import (
	"github.com/foomo/htpasswd"
	log "github.com/sirupsen/logrus"
)

// ModifyAccess grants or revoke usage access for users on a specific repository.
func (svn *SVNMan) ModifyAccess(repoID string, mods ModifyAccess, logFields log.Fields) error {
	logger := log.WithFields(logFields).WithFields(log.Fields{
		"grant_count":  len(mods.Grant),
		"revoke_count": len(mods.Revoke),
	})

	logger.Info("modifying repository access")
	htpasswd.ParseHtpasswdFile("htpasswd")
	return ErrNotImplemented
}
