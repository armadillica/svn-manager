package httphandler

import (
	check "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

/* This doesn't actually test the happy flow, since the mocked HTTP stuff
 * doesn't implement http.Hijacker. */

type HTTPHandlerTestSuite struct {
	fss *APIHandler
}

var _ = check.Suite(&HTTPHandlerTestSuite{})

func (s *HTTPHandlerTestSuite) SetUpTest(c *check.C) {
	httpmock.Activate()
	s.fss = CreateHTTPHandler(nil)
}

func (s *HTTPHandlerTestSuite) TearDownTest(c *check.C) {
	httpmock.DeactivateAndReset()
}
