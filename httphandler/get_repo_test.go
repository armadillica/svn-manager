package httphandler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"

	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

func (s *HTTPHandlerTestSuite) getRepo(c *check.C, repoID string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/unittests/repo/"+repoID, nil)
	req.Header.Set("Content-Type", "application/json")

	respRec := httptest.NewRecorder()
	s.route.ServeHTTP(respRec, req)

	return respRec
}

func (s *HTTPHandlerTestSuite) TestGetRepo(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	mockSVN.EXPECT().GetUsernames("1234").Times(1).Return([]string{}, nil)
	mockSVN.EXPECT().GetUsernames("1234").Times(1).Return([]string{"mysterioususer", "someone.else"}, nil)
	mockSVN.EXPECT().GetUsernames("1234").Times(1).Return(nil, errors.New("test error"))

	resp := RepoDescription{}
	respRec := s.getRepo(c, "1234")
	parseJSON(c, respRec, http.StatusOK, &resp)
	assert.Equal(c, "1234", resp.RepoID)
	assert.Equal(c, []string{}, resp.Access)

	resp = RepoDescription{}
	respRec = s.getRepo(c, "1234")
	parseJSON(c, respRec, http.StatusOK, &resp)
	sort.Strings(resp.Access)
	assert.Equal(c, "1234", resp.RepoID)
	assert.Equal(c, []string{"mysterioususer", "someone.else"}, resp.Access)

	respRec = s.getRepo(c, "1234")
	assert.Equal(c, http.StatusInternalServerError, respRec.Code)
}
