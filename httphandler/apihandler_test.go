package httphandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/armadillica/svn-manager/svnman"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type HTTPHandlerTestSuite struct {
	api *APIHandler
}

var _ = check.Suite(&HTTPHandlerTestSuite{})

func (s *HTTPHandlerTestSuite) SetUpTest(c *check.C) {
	s.api = CreateHTTPHandler(nil)
}

func (s *HTTPHandlerTestSuite) TearDownTest(c *check.C) {
}

func (s *HTTPHandlerTestSuite) TestCreateRepo(c *check.C) {
	mockCtrl := gomock.NewController(c)
	mockSVN := svnman.NewMockManager(mockCtrl) // mocked Manager, not a mock-manager.
	s.api.svn = mockSVN
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:              "4444",
		AuthenticationRealm: "quoted \"strings\" should be válide",
		ProjectID:           "97123333214",
		Creator:             "creator <email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(1)

	body, err := json.Marshal(&repoInfo)
	assert.Nil(c, err, "marshalling failed")
	req, _ := http.NewRequest("POST", "/api/repo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	respRec := httptest.NewRecorder()
	s.api.createRepo(respRec, req)

	assert.Equal(c, 201, respRec.Code)
	assert.Equal(c, "/api/repo/4444", respRec.Header().Get("Location"))
}

// TODO(sybren): test with invalid RepoID, ProjectID, and other values for leakage.
