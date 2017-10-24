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

func (s *ValidationTestSuite) TestValidAuthRealmUnhappy(t *check.C) {
	assert.False(t, ValidAuthRealm(""))
	assert.False(t, ValidAuthRealm("1"))
	assert.False(t, ValidAuthRealm("12"))
	assert.False(t, ValidAuthRealm("123"))
	assert.False(t, ValidAuthRealm("nonbreaking\u00a0space"))
	assert.False(t, ValidAuthRealm("newline\nhere"))
}

func (s *ValidationTestSuite) TestValidAuthRealmHappy(t *check.C) {
	assert.True(t, ValidAuthRealm("peanuts"))
	assert.True(t, ValidAuthRealm("1234"))
	assert.True(t, ValidAuthRealm("a-asdc_D"))
	assert.True(t, ValidAuthRealm("über"))
	assert.True(t, ValidAuthRealm("аррӏе"))
	assert.True(t, ValidAuthRealm("with regular space"))
	assert.True(t, ValidAuthRealm(`realm with "quotes"`))
}
