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

func TestWrite(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		expectPrevious string
		input          string
		expectLatest   string
	}{
		{
			expectPrevious: "",
			input:          "1",
			expectLatest:   "1",
		},
		{
			expectPrevious: "1",
			input:          "2",
			expectLatest:   "12",
		},
		{
			expectPrevious: "12",
			input:          "3",
			expectLatest:   "123",
		},
	}

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer Clean(tempFile)

	for i, test := range testCases {
		err := tempFile.Write(test.input)
		if err != nil {
			t.Fatalf("#%d: %s", i+1, err)
		}
		if string(tempFile.contents[previous]) != test.expectPrevious {
			t.Fatalf("#%d: invalid contents[previous]: expect=%v, but got=%v",
				i+1, test.expectPrevious, string(tempFile.contents[previous]))
		}
		if string(tempFile.contents[latest]) != test.expectLatest {
			t.Fatalf("#%d: invalid contents[latest]: expect=%v, but got=%v",
				i+1, test.expectLatest, string(tempFile.contents[latest]))
		}
	}
}

type mockApp struct{}

func (m *mockApp) run(filePath string) error {
	return nil
}

func TestOpenWith(t *testing.T) {
	t.Parallel()

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer Clean(tempFile)

	editor := &Editor{&mockApp{}}
	tempFile.OpenWith(editor)
	if err != nil {
		t.Fatal(err)
	}
}
