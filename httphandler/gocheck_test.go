/**
 * Common test functionality, and integration with GoCheck.
 */
package httphandler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	check "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
// You only need one of these per package, or tests will run multiple times.
func TestWithGocheck(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	check.TestingT(t)
}

func parseJSON(c *check.C, respRec *httptest.ResponseRecorder, expectedStatus int, parsed interface{}) {
	assert.Equal(c, expectedStatus, respRec.Code)
	headers := respRec.Header()
	assert.Equal(c, "application/json", headers.Get("Content-Type"))

	decoder := json.NewDecoder(respRec.Body)
	if err := decoder.Decode(&parsed); err != nil {
		c.Fatalf("Unable to decode JSON: %s", err)
	}
}
