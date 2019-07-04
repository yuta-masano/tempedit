package tempedit

import (
	"os"
	"testing"
)

func TestNewTempFileAndClean(t *testing.T) {
	t.Parallel()

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func(tempFile *TempFile) {
		Clean(tempFile)
		_, err = os.Stat(tempFile.Name())
		if err == nil {
			t.Fatalf("%s should be deleted but exists", tempFile.Name())
		}
	}(tempFile)

	_, err = os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("%s should be created but not found", tempFile.Name())
	}

	if string(tempFile.contents[previous]) != "" && string(tempFile.contents[latest]) != "" {
		t.Fatalf("both of internal contents should be empty string but previous: %s, latest: %s",
			string(tempFile.contents[previous]), string(tempFile.contents[latest]))
	}
}
