package httphandler

import (
	"path/filepath"
	"regexp"

	"github.com/armadillica/svn-manager/filelocator"
	"github.com/xeipuuv/gojsonschema"
)

var (
	validRepoRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-]+[a-zA-Z0-9]$`)
)

// ValidRepoID returns true iff the repoID is safe to use as SVN repository name/path/ID.
func ValidRepoID(repoID string) bool {
	if len(repoID) < 4 {
		return false
	}
	return validRepoRegexp.MatchString(repoID)
}

// ValidRequest validates the given document against the given schema.
func validRequest(schemaName string, document interface{}) (*gojsonschema.Result, error) {
	filename, err := filelocator.FindFile(filepath.Join("json_schemas", schemaName+".json"))
	if err != nil {
		return nil, err
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + filename)
	documentLoader := gojsonschema.NewGoLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, err
	}
	return result, nil
}
