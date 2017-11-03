package svnman

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

func (s *SVNManTestSuite) TestDeleteHappy(t *check.C) {
	logFields := log.Fields{"in": "unittest"}

	repoInfo := CreateRepo{
		RepoID:    "my-repo-id",
		ProjectID: "59eefa9cf488554678cae036",
		Creator:   "dr. Stüvel <sybren@blender.studio>",
	}
	if err := s.svn.CreateRepo(repoInfo, logFields); err != nil {
		t.Fatalf("Unable to create repo: %s", err)
	}
	// Any restarts queued by CreateRepo are irrelevant to this test.
	s.mr = mockRestarter{}

	if err := s.svn.DeleteRepo("my-repo-id", logFields); err != nil {
		t.Fatalf("unexpected error deleting repo: %s", err)
	}

	rp := s.svn.repoPath("my-repo-id")
	if _, err := os.Stat(rp); err == nil {
		t.Errorf("repo path %q should not exist after deletion", rp)
	}

	// the repository should be moved into the attic.
	glob := filepath.Join(s.svn.repoRoot, "attic", "my", "my-repo-id-2*")
	found, err := filepath.Glob(glob)
	if err != nil {
		t.Fatalf("error globbing %s: %s", glob, err)
	}
	assert.Equal(t, 1, len(found), "the repository should be moved into the attic")

	// same for the Apache configuration file.
	glob = filepath.Join(s.svn.apacheConfigDir, "attic", "my", "svn-my-repo-id.conf-2*")
	found, err = filepath.Glob(glob)
	if err != nil {
		t.Fatalf("error globbing %s: %s", glob, err)
	}
	assert.Equal(t, 1, len(found), "the Apache config file should be moved into the attic")

	assert.True(t, s.mr.restartCalled, "an Apache restart should have been queued")
}

func (s *SVNManTestSuite) TestDeleteNonExistantRepoHappy(t *check.C) {
	logFields := log.Fields{"in": "unittest"}
	err := s.svn.DeleteRepo("my-repo-id", logFields)
	assert.Nil(t, err)
	assert.True(t, s.mr.restartCalled, "an Apache restart should have been queued")
}
