package svnman

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

func (s *SVNManTestSuite) loadHtpasswd(t *check.C, repoID string) []string {
	repo := filepath.Join(s.svn.repoRoot, repoID[0:2], repoID)
	passwdfile := filepath.Join(repo, "htpasswd")
	htpasswdBytes, err := ioutil.ReadFile(passwdfile)
	if err != nil {
		t.Fatalf("Unable to open %s: %s", passwdfile, err)
	}
	htpasswd := string(htpasswdBytes)
	lines := strings.Split(strings.TrimSpace(htpasswd), "\n")
	return lines
}

func (s *SVNManTestSuite) TestModifyAccessHappy(t *check.C) {
	logFields := log.Fields{"in": "unittest"}

	repoInfo := CreateRepo{
		RepoID:    "1234",
		ProjectID: "59eefa9cf488554678cae036",
		Creator:   "dr. Stüvel <sybren@blender.studio>",
	}
	if err := s.svn.CreateRepo(repoInfo, logFields); err != nil {
		t.Fatalf("Unable to create repo: %s", err)
	}

	// Grant access to one user.
	if err := s.svn.ModifyAccess(repoInfo.RepoID, ModifyAccess{
		Grant: []ModifyAccessGrantEntry{
			ModifyAccessGrantEntry{"testkees", "$2y$05$cWVQLHS58K7fIKjz3tU52eBI2sxbE3KdAfZN0CJN"},
		},
	}, logFields); err != nil {
		t.Fatalf("Unable to modify access: %s", err)
	}

	lines := s.loadHtpasswd(t, repoInfo.RepoID)
	assert.Equal(t, 1, len(lines), "strange line count, file content: %s", strings.Join(lines, `\\`))
	oneline := strings.SplitN(lines[0], ":", 2)
	assert.Equal(t, "testkees", oneline[0])
	assert.Equal(t, "$2y$05$cWVQLHS58K7fIKjz3tU52eBI2sxbE3KdAfZN0CJN", oneline[1])

	// Modify password of one user, and grant access to a new one.
	if err := s.svn.ModifyAccess(repoInfo.RepoID, ModifyAccess{
		Grant: []ModifyAccessGrantEntry{
			ModifyAccessGrantEntry{"testkees", "$2y$05$cWZN0CJN"},
			ModifyAccessGrantEntry{"anotherone", "$2y$05$cW---ZN0CJN"},
		},
	}, logFields); err != nil {
		t.Fatalf("Unable to re-modify access: %s", err)
	}

	lines = s.loadHtpasswd(t, repoInfo.RepoID)
	assert.Equal(t, 2, len(lines))

	// Those two lines can be in any order.
	found := map[string]string{}
	for _, line := range lines {
		words := strings.SplitN(line, ":", 2)
		found[words[0]] = words[1]
	}
	assert.Equal(t, map[string]string{
		"testkees":   "$2y$05$cWZN0CJN",
		"anotherone": "$2y$05$cW---ZN0CJN",
	}, found)

	// Revoke access from one existing and one non-existing user.
	if err := s.svn.ModifyAccess(repoInfo.RepoID, ModifyAccess{
		Revoke: []string{"testkees", "nonexisting"},
	}, logFields); err != nil {
		t.Fatalf("Unable to re-modify access: %s", err)
	}

	lines = s.loadHtpasswd(t, repoInfo.RepoID)
	assert.Equal(t, 1, len(lines))
	oneline = strings.SplitN(lines[0], ":", 2)
	assert.Equal(t, "anotherone", oneline[0])
	assert.Equal(t, "$2y$05$cW---ZN0CJN", oneline[1])
}
