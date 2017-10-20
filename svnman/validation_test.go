package svnman

import (
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type ValidationTestSuite struct {
}

var _ = check.Suite(&ValidationTestSuite{})

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
