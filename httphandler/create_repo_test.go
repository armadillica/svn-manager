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

func (s *HTTPHandlerTestSuite) createRepo(c *check.C, repoInfo svnman.CreateRepo) *httptest.ResponseRecorder {
	body, err := json.Marshal(repoInfo)
	assert.Nil(c, err, "marshalling failed")

	req, _ := http.NewRequest("POST", "/unittests/repo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	respRec := httptest.NewRecorder()
	s.route.ServeHTTP(respRec, req)

	return respRec
}

func (s *HTTPHandlerTestSuite) TestCreateRepo(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:    "4444",
		ProjectID: "8afae1eb1d171833df73416b",
		Creator:   "creator <email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(1)

	respRec := s.createRepo(c, repoInfo)

	assert.Equal(c, 201, respRec.Code)
	assert.Equal(c, "/unittests/repo/4444", respRec.Header().Get("Location"))
}

func (s *HTTPHandlerTestSuite) TestCreateRepoBadRepoID(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:    "in valid",
		ProjectID: "8afae1eb1d171833df73416b",
		Creator:   "creator <email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(0)

	respRec := s.createRepo(c, repoInfo)
	assert.Equal(c, http.StatusBadRequest, respRec.Code)
}

func (s *HTTPHandlerTestSuite) TestCreateRepoAlreadyExists(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:    "alreadyexists",
		ProjectID: "8afae1eb1d171833df73416b",
		Creator:   "creator <email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(1).Return(svnman.ErrAlreadyExists)

	respRec := s.createRepo(c, repoInfo)
	assert.Equal(c, http.StatusConflict, respRec.Code)
}

func (s *HTTPHandlerTestSuite) TestCreateRepoBadProjectID(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:    "valid",
		ProjectID: "8afae1eb1d171\n833df73416b",
		Creator:   "creator <email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(0)

	respRec := s.createRepo(c, repoInfo)
	assert.Equal(c, http.StatusBadRequest, respRec.Code)
}

func (s *HTTPHandlerTestSuite) TestCreateRepoBadCreator(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	repoInfo := svnman.CreateRepo{
		RepoID:    "valid",
		ProjectID: "8afae1eb1d171833df73416b",
		Creator:   "creator\n<email@example.com>",
	}

	mockSVN.EXPECT().CreateRepo(repoInfo, gomock.Any()).Times(0)

	respRec := s.createRepo(c, repoInfo)
	assert.Equal(c, http.StatusBadRequest, respRec.Code)
}
