package httphandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// decodeJSON decodes JSON from an io.Reader, and writes a Bad Request status if it fails.
func decodeJSON(w http.ResponseWriter, r io.Reader, document interface{}, logFields log.Fields) error {
	dec := json.NewDecoder(r)

	if err := dec.Decode(document); err != nil {
		log.WithFields(logFields).WithError(err).Warning("unable to decode JSON")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unable to decode JSON: %s\n", err)
		return err
	}

	return nil
}
