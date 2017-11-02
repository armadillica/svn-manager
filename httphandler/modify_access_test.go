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

func (s *HTTPHandlerTestSuite) modifyAccess(c *check.C, repoID string, payload svnman.ModifyAccess) *httptest.ResponseRecorder {
	body, err := json.Marshal(payload)
	assert.Nil(c, err, "marshalling failed")

	req, _ := http.NewRequest("POST", "/unittests/repo/"+repoID+"/access", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	respRec := httptest.NewRecorder()
	s.route.ServeHTTP(respRec, req)

	return respRec
}

func (s *HTTPHandlerTestSuite) TestModifyAccessHappy(c *check.C) {
	mockCtrl, mockSVN := s.mockSVN(c)
	defer mockCtrl.Finish()

	payload := svnman.ModifyAccess{
		Grant: []svnman.ModifyAccessGrantEntry{svnman.ModifyAccessGrantEntry{
			Username: "mysterioususer",
			Password: "$2y$10$abcdef",
		}},
	}

	mockSVN.EXPECT().ModifyAccess("1234", payload, gomock.Any()).Times(1)

	respRec := s.modifyAccess(c, "1234", payload)
	assert.Equal(c, http.StatusOK, respRec.Code)
}
