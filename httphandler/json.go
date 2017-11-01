package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// ErrBadContentType indicates that the content type of an HTTP request wasn't as expected.
	ErrBadContentType = errors.New("bad content type")

	// ErrValidationFailed indicates that JSON validation failed; details are logged & returned via HTTP.
	ErrValidationFailed = errors.New("JSON validation failed")
)

// decodeJSON decodes JSON from a HTTP request, and writes a Bad Request status if it fails.
func decodeJSON(w http.ResponseWriter, r *http.Request, document interface{}, schemaName string, logFields log.Fields) error {
	ct := r.Header.Get("Content-Type")
	if ct != "application/json" {
		log.WithFields(logFields).WithField("content_type", ct).Warning("expected JSON, got different content")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "expected application/json content type")
		return ErrBadContentType
	}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(document); err != nil {
		log.WithFields(logFields).WithError(err).Warning("unable to decode JSON")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to decode JSON: %s\n", err)
		return ErrValidationFailed
	}

	result, err := validRequest(schemaName, document)
	if err != nil || !result.Valid() {
		writeValidationError(result, err, w, logFields)
		return ErrValidationFailed
	}

	return nil
}

func writeValidationError(result *gojsonschema.Result, err error, w http.ResponseWriter, logFields log.Fields) {
	logger := log.WithFields(logFields).WithError(err)

	if err != nil {
		logger.Warning("received JSON that was unvalidatable")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "unable to validate your JSON")
		return
	}
	if result.Valid() {
		logger.Error("writeValidationError() called for valid JSON")
		return
	}

	verrors := result.Errors()
	errors := make([]string, len(verrors))
	for idx, verr := range verrors {
		errors[idx] = verr.String()
	}
	logger.WithField(log.ErrorKey, strings.Join(errors, "|")).Warning("received invalid JSON")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "unable to validate your JSON:\n")
	for _, errstr := range errors {
		fmt.Fprintf(w, "  - %s\n", errstr)
	}
}
