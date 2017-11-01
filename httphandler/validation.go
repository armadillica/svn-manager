package httphandler

import (
	"go/build"
	"os"
	"path/filepath"
	"regexp"

	"github.com/armadillica/svn-manager/svnman"
	"github.com/kardianos/osext"
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

// Returns the filename as an absolute path.
// Relative paths are interpreted relative to the executable.
func relativeToExecutable(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	}

	exedirname, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	return filepath.Join(exedirname, filename), nil
}

// Returns the filename as an absolute path.
// Relative paths are interpreted relative to the executable.
func relativeToCwd(filename string) (string, error) {
	if filepath.IsAbs(filename) {
		return filename, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, filename), nil
}

func findSchema(schemaName string) (string, error) {
	var schemaFilename string

	// Search relative to the executable, for production.
	schemaFilename, err := relativeToExecutable(schemaName)
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(schemaFilename); err == nil {
		// Found it!
		return schemaFilename, nil
	}

	// Search in the current directory, for development/debug runs.
	schemaFilename, err = relativeToCwd(schemaName)
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(schemaFilename); err == nil {
		// Found it!
		return schemaFilename, nil
	}

	// Search in the build directory, for unit tests.
	p, err := build.Default.Import("github.com/armadillica/svn-manager", "", build.FindOnly)
	if err != nil {
		return "", err
	}
	schemaFilename = filepath.Join(p.Dir, schemaName)
	if _, err = os.Stat(schemaFilename); err == nil {
		// Found it!
		return schemaFilename, nil
	}

	return "", err
}

// ValidRequest validates the given document against the given schema.
func validRequest(schemaName string, document interface{}) (*gojsonschema.Result, error) {
	filename, err := findSchema(filepath.Join("json_schemas", schemaName+".json"))
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

// ValidCreateRepo validates the CreateRepo document, returning validation results.
func ValidCreateRepo(document *svnman.CreateRepo) (*gojsonschema.Result, error) {
	return validRequest("create_repo", document)
}
