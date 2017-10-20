package httphandler

import (
	check "gopkg.in/check.v1"
	"gopkg.in/jarcoal/httpmock.v1"
)

type HTTPHandlerTestSuite struct {
	api *APIHandler
}

var _ = check.Suite(&HTTPHandlerTestSuite{})

func (s *HTTPHandlerTestSuite) SetUpTest(c *check.C) {
	httpmock.Activate()
	s.api = CreateHTTPHandler(nil)
}

func (s *HTTPHandlerTestSuite) TearDownTest(c *check.C) {
	httpmock.DeactivateAndReset()
}

func (s *HTTPHandlerTestSuite) TestCreateRepo(c *check.C) {

}
