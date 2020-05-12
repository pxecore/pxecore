package locator

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/errors"
	"io"
	"os"
	"strings"
)

// StaticFile returns the same locator for all mac addresses.
type StaticFile struct {
	BaseDir  string
	BasePath string
}

// NewStaticFile instantiates a new SingleIPXEScript with it's read-only attributes
func NewStaticFile(baseDir string, basePath string) *StaticFile {
	return &StaticFile{
		BaseDir:  baseDir,
		BasePath: basePath,
	}
}

// Lookup returns the locator for the provided path.
// See gitlab.com/pliego/pxecore/pkg/tftp/FileLocator
func (s StaticFile) Lookup(path string) (io.Reader, error) {
	if s.BasePath == "" || s.BaseDir == "" || !strings.HasPrefix(path, s.BasePath) {
		return nil, &errors.Error{Code: errors.ENotFound,
			Msg: "[tftp.locator] file not found."}
	}
	file, err := os.Open(fmt.Sprint(s.BaseDir, path))
	if err != nil {
		return nil, &errors.Error{Code: errors.ENotFound,
			Msg: "[tftp.locator] error retrieving file.", Err: err}
	}
	return file, nil
}
