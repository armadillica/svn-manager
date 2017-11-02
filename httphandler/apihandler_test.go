package httphandler

import (
	"github.com/armadillica/svn-manager/svnman"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	check "gopkg.in/check.v1"
)

type HTTPHandlerTestSuite struct {
	api   *APIHandler
	route *mux.Router
}

var _ = check.Suite(&HTTPHandlerTestSuite{})

func (s *HTTPHandlerTestSuite) SetUpTest(c *check.C) {
	s.route = mux.NewRouter()
	s.api = CreateAPIHandler(nil)
	s.api.AddRoutes(s.route.PathPrefix("/unittests").Subrouter())
}

func (s *HTTPHandlerTestSuite) TearDownTest(c *check.C) {
}

func (s *HTTPHandlerTestSuite) mockSVN(c *check.C) (*gomock.Controller, *svnman.MockManager) {
	mockCtrl := gomock.NewController(c)
	mockSVN := svnman.NewMockManager(mockCtrl) // mocked Manager, not a mock-manager.
	s.api.svn = mockSVN

	return mockCtrl, mockSVN
}
