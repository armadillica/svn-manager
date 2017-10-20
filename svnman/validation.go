package svnman

import (
	"regexp"
)

var validRepoRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]+[a-zA-Z0-9]$`)

// ValidRepoID returns true iff the repoID is safe to use as SVN repository name/path/ID.
func ValidRepoID(repoID string) bool {
	if len(repoID) < 4 {
		return false
	}
	return validRepoRegexp.MatchString(repoID)
}
