package svnman

import "errors"

// SVNMan provides SVN management operations.
type SVNMan struct {
}

// CreateRepo creates a repository and Apache location directive.
func (svn *SVNMan) CreateRepo(repoInfo CreateRepo) error {
	return errors.New("not implemented")
}
