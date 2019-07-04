package tempedit

import (
	"io/ioutil"
	"os"
)

const (
	// these are used for the index of tempFile.conntents slice.
	previous int = iota
	latest
)

// TempFile is a temporary file object to be written.
type TempFile struct {
	*tempFile
}

type tempFile struct {
	*os.File
	// In order to detect changes the contents before and after, the change are
	// held.
	contents [2][]byte
}

// NewTempFile creates a new temporary file.
// tmpDir and prefix are same as the args of ioutil.TempFile function.
func NewTempFile(tmpDir, prefix string) (*TempFile, error) {
	f := &TempFile{&tempFile{}}
	var err error
	f.File, err = ioutil.TempFile(tmpDir, prefix)
	if err != nil {
		return nil, err
	}

	f.contents[previous], f.contents[latest] = []byte(""), []byte("")
	return f, nil
}

// Clean is used after NewTempFile.
// You have responsible to defer this function after NewTempFile.
func Clean(tempFiles ...*TempFile) {
	for _, f := range tempFiles {
		if err := f.Close(); err != nil {
			panic(err)
		}
		if err := os.Remove(f.Name()); err != nil {
			panic(err)
		}
	}
}
