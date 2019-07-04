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

func TestWriteTEmplate(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		expectPrevious string
		src            string
		data           interface{}
		expectLatest   string
	}{
		{
			expectPrevious: "",
			src:            "Test {{ .Num }}",
			data:           struct{ Num string }{Num: "1"},
			expectLatest:   "Test 1",
		},
		{
			expectPrevious: "Test 1",
			src:            "Test {{ .Num }} and `" + `!"#$%&'()=~|{+*}<>?_\^-@[]:;\/.,`,
			data:           struct{ Num string }{Num: "2"},
			expectLatest:   "Test 1Test 2 and `" + `!"#$%&'()=~|{+*}<>?_\^-@[]:;\/.,`,
		},
		{
			expectPrevious: "Test 1Test 2 and `" + `!"#$%&'()=~|{+*}<>?_\^-@[]:;\/.,`,
			src:            "Test {{ .Str }}",
			data:           struct{ Str string }{Str: "3"},
			expectLatest:   "Test 1Test 2 and `" + `!"#$%&'()=~|{+*}<>?_\^-@[]:;\/.,` + "Test 3",
		},
	}

	tempFile, err := NewTempFile("", "")
	defer Clean(tempFile)
	if err != nil {
		t.Fatal(err)
	}

	for i, test := range testCases {
		err := tempFile.WriteTemplate(test.src, test.data)
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
