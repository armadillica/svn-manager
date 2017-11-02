// Package filelocator helps to find files relative to CWD, exectuable path, and our source path.
package filelocator

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
)

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

// FindFile returns the absolute path of the named file.
func FindFile(filename string) (string, error) {
	var schemaFilename string

	// Search relative to the executable, for production.
	schemaFilename, err := relativeToExecutable(filename)
	if err != nil {
		return "", err
	}
	if _, err = os.Stat(schemaFilename); err == nil {
		// Found it!
		return schemaFilename, nil
	}

	// Search in the current directory, for development/debug runs.
	schemaFilename, err = relativeToCwd(filename)
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
	schemaFilename = filepath.Join(p.Dir, filename)
	if _, err = os.Stat(schemaFilename); err == nil {
		// Found it!
		return schemaFilename, nil
	}

	return "", err
}
