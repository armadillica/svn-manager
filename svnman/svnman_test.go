package svnman

import (
	"io/ioutil"
	"os"

	check "gopkg.in/check.v1"
)

type mockRestarter struct {
	restartCalled bool
}

func (mr *mockRestarter) QueueRestart() {
	mr.restartCalled = true
}

type SVNManTestSuite struct {
	svn *SVNMan
	mr  mockRestarter
}

var _ = check.Suite(&SVNManTestSuite{})

func mustTempDir(dir, prefix string) string {
	tempdir, err := ioutil.TempDir(dir, prefix)
	if err != nil {
		panic(err)
	}
	return tempdir
}

func (s *SVNManTestSuite) SetUpTest(t *check.C) {
	s.svn = &SVNMan{
		restarter:       &s.mr,
		repoRoot:        mustTempDir("", "reporoot"),
		apacheConfigDir: mustTempDir("", "apache"),
		appName:         "SVNMan unit test",
		appVersion:      "0.1.2.3-beta5-sub3",
	}
}

func (s *SVNManTestSuite) TearDownTest(t *check.C) {
	if err := os.RemoveAll(s.svn.repoRoot); err != nil {
		t.Fatal("unable to remove repoRoot", s.svn.repoRoot, err)
	}
	if err := os.RemoveAll(s.svn.apacheConfigDir); err != nil {
		t.Fatal("unable to remove apacheConfigDir", s.svn.apacheConfigDir, err)
	}
}
