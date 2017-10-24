package svnman

import (
	"regexp"
)

var (
	validRepoRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]+[a-zA-Z0-9]$`)
	validAuthRealm  = regexp.MustCompile(`^[\p{L}\d_\- "']+$`)
)

// ValidRepoID returns true iff the repoID is safe to use as SVN repository name/path/ID.
func ValidRepoID(repoID string) bool {
	if len(repoID) < 4 {
		return false
	}
	return validRepoRegexp.MatchString(repoID)
}

// ValidAuthRealm returns true iff the realm is safe to use as Apache authentication realm.
func ValidAuthRealm(realm string) bool {
	if len(realm) < 4 {
		return false
	}
	return validAuthRealm.MatchString(realm)
}
