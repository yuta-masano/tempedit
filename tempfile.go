package tempedit

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"
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

// pushContent must be called every writing method.
// It pushes contents of before / after writing the temporary file to
// tempFile.contents.
func (t *tempFile) pushContent() error {
	_, err := t.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(t)
	if err != nil {
		return err
	}
	t.contents[previous], t.contents[latest] = t.contents[latest], b
	return nil
}

// Write writes a content into the temporary file.
func (t *tempFile) Write(content string) error {
	_, err := t.File.Write([]byte(content))
	if err != nil {
		return err
	}
	return t.pushContent()
}

// OpenWith opens the temporary file with external editor.
func (t *tempFile) OpenWith(editor *Editor) error {
	err := editor.run(t.Name())
	if err != nil {
		return err
	}
	return t.pushContent()
}

// WriteTemplate is a wrapper method of text/template.
// src is string to be parsed and data is objedt to be applied parsed template.
func (t *tempFile) WriteTemplate(src string, data interface{}) error {
	tpl, err := template.New("").Parse(src)
	if err != nil {
		return err
	}
	err = tpl.Execute(t.File, data)
	if err != nil {
		return err
	}
	return t.pushContent()
}
