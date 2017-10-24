package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ErrBadContentType indicates that the content type of an HTTP request wasn't as expected.
var ErrBadContentType = errors.New("bad content type")

// decodeJSON decodes JSON from a HTTP request, and writes a Bad Request status if it fails.
func decodeJSON(w http.ResponseWriter, r *http.Request, document interface{}, logFields log.Fields) error {
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
		return err
	}

	return nil
}
