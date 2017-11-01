package httphandler

import (
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type ValidationTestSuite struct {
}

var _ = check.Suite(&ValidationTestSuite{})

func (s *ValidationTestSuite) assertValidJSON(t *check.C, schemaName string, doc interface{}) {
	result, err := validRequest(schemaName, doc)
	if err != nil {
		t.Errorf("validation of request failed completely: %s", err)
		return
	}

	assert.True(t, result.Valid(), "unexpectedly invalid document: %s", result.Errors())
}

func (s *ValidationTestSuite) assertInvalidJSON(t *check.C, schemaName string, doc interface{}) {
	result, err := validRequest(schemaName, doc)
	if err != nil {
		t.Errorf("validation of request failed completely: %s", err)
		return
	}
	assert.False(t, result.Valid(), "validation unexpectedly succeeded with doc: %s", doc)
}
