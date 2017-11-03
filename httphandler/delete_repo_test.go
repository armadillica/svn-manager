package httphandler

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/armadillica/svn-manager/svnman"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

func (s *HTTPHandlerTestSuite) deleteRepo(c *check.C, repoID string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", "/unittests/repo/"+repoID, nil)
	respRec := httptest.NewRecorder()
	s.route.ServeHTTP(respRec, req)
	return respRec
}

func (s *HTTPHandlerTestSuite) TestDeleteRepo(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	mockSVN.EXPECT().DeleteRepo("1234", gomock.Any()).Times(1).Return(nil)
	mockSVN.EXPECT().DeleteRepo("12345", gomock.Any()).Times(1).Return(svnman.ErrNotFound)
	mockSVN.EXPECT().DeleteRepo("123456", gomock.Any()).Times(1).Return(errors.New("something unexpected"))

	respRec := s.deleteRepo(c, "1234")
	assert.Equal(c, http.StatusNoContent, respRec.Code)

	respRec = s.deleteRepo(c, "12345")
	assert.Equal(c, http.StatusNotFound, respRec.Code)

	respRec = s.deleteRepo(c, "123456")
	assert.Equal(c, http.StatusInternalServerError, respRec.Code)
}
