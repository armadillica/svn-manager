package httphandler

import (
	"github.com/armadillica/svn-manager/svnman"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

var validCreateRepo = svnman.CreateRepo{
	RepoID:    "peanuts",
	ProjectID: "8afae1eb1d171833df73416b",
	Creator:   "Mrs. Unícøde <some@example.com>",
}

func (s *ValidationTestSuite) TestValidRepoIDUnhappy(t *check.C) {
	assert.False(t, ValidRepoID(""))
	assert.False(t, ValidRepoID("1"))
	assert.False(t, ValidRepoID("12"))
	assert.False(t, ValidRepoID("123"))
	assert.False(t, ValidRepoID("-asdc_"))
	assert.False(t, ValidRepoID("über"))
	assert.False(t, ValidRepoID("аррӏе"))
	assert.False(t, ValidRepoID("with regular space"))
	assert.False(t, ValidRepoID("nonbreaking\u00a0space"))
}

func (s *ValidationTestSuite) TestValidRepoIDHappy(t *check.C) {
	assert.True(t, ValidRepoID("peanuts"))
	assert.True(t, ValidRepoID("1234"))
	assert.True(t, ValidRepoID("a-asdc_D"))
}

func (s *ValidationTestSuite) TestValidRepoIDUnhappyJSON(t *check.C) {
	s.testInvalidRepoID(t, "")
	s.testInvalidRepoID(t, "1")
	s.testInvalidRepoID(t, "12")
	s.testInvalidRepoID(t, "123")
	s.testInvalidRepoID(t, "-asdc_")
	s.testInvalidRepoID(t, "über")
	s.testInvalidRepoID(t, "аррӏе")
	s.testInvalidRepoID(t, "with regular space")
	s.testInvalidRepoID(t, "nonbreaking\u00a0space")
}

func (s *ValidationTestSuite) TestValidRepoIDHappyJSON(t *check.C) {
	s.testValidRepoID(t, "peanuts")
	s.testValidRepoID(t, "1234")
	s.testValidRepoID(t, "a-asdc_D")
}

func (s *ValidationTestSuite) TestValidProjectIDUnhappy(t *check.C) {
	s.testInvalidProjectID(t, "")
	s.testInvalidProjectID(t, "1")
	s.testInvalidProjectID(t, "12")
	s.testInvalidProjectID(t, "123")
	s.testInvalidProjectID(t, "aaaaaaaaaaaaaaaaaaaaaaaaa")
	s.testInvalidProjectID(t, "a1dac95688bäd8b045a073b0")
	s.testInvalidProjectID(t, "nonbreaking\u00a0space")
	s.testInvalidProjectID(t, "newline\nhere")
}

func (s *ValidationTestSuite) TestValidProjectIDHappy(t *check.C) {
	s.testValidProjectID(t, "8afae1eb1d171833df73416b")
	s.testValidProjectID(t, "aaaaaaaaaaaaaaaaaaaaaaaa")
}

func (s *ValidationTestSuite) TestValidCreatorHappy(t *check.C) {
	s.testValidCreator(t, validCreateRepo.Creator)
	s.testValidCreator(t, "8afae1eb1d171833df73416@asdd130j")
	s.testValidCreator(t, "Mr. ASCII <some@example.com>")
	s.testValidCreator(t, "b46ff9af18335ee3367411949e583163bb37ac6ac63ee03cd1e08a1be27a853e2c0fa8734d1d3ae89371b47033a2bbc818d6893069209f92466bde28af178883")
}

func (s *ValidationTestSuite) TestValidCreatorUnhappy(t *check.C) {
	s.testInvalidCreator(t, "")
	s.testInvalidCreator(t, "123")
	s.testInvalidCreator(t, "nonbreaking\u00a0space")
	s.testInvalidCreator(t, "Mr. New\nline <some@example.com>")
	s.testInvalidCreator(t, "b46ff9af18335ee3367411949e583163bb37ac6ac63ee03cd1e08a1be27a853e2c0fa8734d1d3ae89371b47033a2bbc818d6893069209f92466bde28af178883d5")
}

func (s *ValidationTestSuite) testValidRepoID(t *check.C, repoID string) {
	doc := validCreateRepo
	doc.RepoID = repoID
	s.assertValidJSON(t, "create_repo", &doc)
}

func (s *ValidationTestSuite) testInvalidRepoID(t *check.C, repoID string) {
	doc := validCreateRepo
	doc.RepoID = repoID
	s.assertInvalidJSON(t, "create_repo", &doc)
}

func (s *ValidationTestSuite) testValidCreator(t *check.C, creator string) {
	doc := validCreateRepo
	doc.Creator = creator
	s.assertValidJSON(t, "create_repo", &doc)
}

func (s *ValidationTestSuite) testInvalidCreator(t *check.C, creator string) {
	doc := validCreateRepo
	doc.Creator = creator
	s.assertInvalidJSON(t, "create_repo", &doc)
}

func (s *ValidationTestSuite) testValidProjectID(t *check.C, projectID string) {
	doc := validCreateRepo
	doc.ProjectID = projectID
	s.assertValidJSON(t, "create_repo", &doc)
}

func (s *ValidationTestSuite) testInvalidProjectID(t *check.C, projectID string) {
	doc := validCreateRepo
	doc.ProjectID = projectID
	s.assertInvalidJSON(t, "create_repo", &doc)
}
