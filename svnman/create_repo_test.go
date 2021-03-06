package svnman

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

func (s *SVNManTestSuite) TestCreateRepoHappy(t *check.C) {
	repoInfo := CreateRepo{
		RepoID:    "1234",
		ProjectID: "59eefa9cf488554678cae036",
		Creator:   "dr. Stüvel <sybren@blender.studio>",
	}

	logFields := log.Fields{"in": "unittest"}
	err := s.svn.CreateRepo(repoInfo, logFields)
	assert.Nil(t, err, "unable to create repo: %s", err)

	// Check repository files.
	repo := filepath.Join(s.svn.repoRoot, "12", repoInfo.RepoID)
	stat, err := os.Stat(repo)
	if err != nil {
		assert.Fail(t, err.Error(), "repo %q should exist", repo)
	} else {
		assert.True(t, stat.IsDir(), "repo %q should be a directory", repo)
	}

	fmtfile := filepath.Join(repo, "format")
	fmtcontent, err := ioutil.ReadFile(fmtfile)
	if err != nil {
		assert.Fail(t, err.Error(), "file %q should exist", fmtfile)
	} else {
		assert.Equal(t, "5", strings.TrimSpace(string(fmtcontent)))
	}

	passwdfile := filepath.Join(repo, "htpasswd")
	stat, err = os.Stat(passwdfile)
	if err != nil {
		assert.Fail(t, err.Error(), "password file %q should exist", passwdfile)
	} else {
		assert.Equal(t, int64(0), stat.Size(), "password file %q should be empty", passwdfile)
	}

	infofile := filepath.Join(repo, "info.yaml")
	infobytes, err := ioutil.ReadFile(infofile)
	if err != nil {
		assert.Fail(t, err.Error(), "info file %q should exist", infofile)
	} else {
		info := string(infobytes)
		assert.Contains(t, info, "dr. Stüvel <sybren@blender.studio>")
		assert.Contains(t, info, "1234")
		assert.Contains(t, info, "59eefa9cf488554678cae036")
	}

	// Check Apache location directive file.
	apache := filepath.Join(s.svn.apacheConfigDir, "12", "svn-"+repoInfo.RepoID+".conf")
	apabytes, err := ioutil.ReadFile(apache)
	if err != nil {
		assert.Fail(t, err.Error(), "file %q should exist", apache)
	} else {
		apa := string(apabytes)
		assert.Contains(t, apa, "/repo/1234")
		assert.Contains(t, apa, repo)
		assert.Contains(t, apa, "59eefa9cf488554678cae036", "Project ID should be mentioned in Apache config file")
		assert.Contains(t, apa, `\"1234\"`, "Auth realm should be quoted properly")
	}

	assert.True(t, s.mr.restartCalled, "Apache restart not requested")
}

func (s *SVNManTestSuite) TestCreateRepoAlreadyExists(t *check.C) {
	repoInfo := CreateRepo{
		RepoID:    "1234",
		ProjectID: "59eefa9cf488554678cae036",
		Creator:   "dr. Stüvel <sybren@blender.studio>",
	}

	logFields := log.Fields{"in": "unittest"}
	err := s.svn.CreateRepo(repoInfo, logFields)
	assert.Nil(t, err, "unable to create repo: %s", err)

	err = s.svn.CreateRepo(repoInfo, logFields)
	assert.Equal(t, ErrAlreadyExists, err)
}
