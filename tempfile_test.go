package tempedit

import (
	"bytes"
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

func TestIsChangedCaseNotChanged(t *testing.T) {
	t.Parallel()

	var err error

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer Clean(tempFile)
	_, err = tempFile.IsChanged()
	if err.Error() != ErrMsgNotChanged {
		t.Fatalf("temporary file initialize must be non-changed status: expect=%s, but got=%s", ErrMsgNotChanged, err)
	}

	var caseNotChanged = []struct {
		expect string
		before string
		after  string
	}{
		{
			expect: ErrMsgNotChanged,
			before: "foo",
			after:  "foo",
		},
		{
			expect: ErrMsgNotChanged,
			before: "foo\n",
			after:  "foo\n",
		},
		{
			expect: ErrMsgNotChanged,
			before: "foo",
			after:  "foo\n",
		},
		{
			expect: ErrMsgNotChanged,
			before: "foo\n",
			after:  "foo",
		},
	}
	for i, test := range caseNotChanged {
		tempFile.contents[previous], tempFile.contents[latest] =
			[]byte(test.before), []byte(test.after)
		_, err = tempFile.IsChanged()
		if err.Error() != ErrMsgNotChanged {
			t.Fatalf("caseNotChanged: #%d: wrong err msg: expect=%s, but got=%s", i+1, ErrMsgNotChanged, err)
		}
	}
}

func TestIsChangedCaseEmpty(t *testing.T) {
	t.Parallel()

	var err error

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer Clean(tempFile)
	_, err = tempFile.IsChanged()
	if err.Error() != ErrMsgNotChanged {
		t.Fatalf("temporary file initialize must be non-changed status: expect=%s, but got=%s", ErrMsgNotChanged, err)
	}

	var caseEmpty = struct {
		expect string
		before string
		after  string
	}{
		expect: ErrMsgEmpty,
		before: "foo",
		after:  "",
	}
	tempFile.contents[previous], tempFile.contents[latest] =
		[]byte(caseEmpty.before), []byte(caseEmpty.after)
	_, err = tempFile.IsChanged()
	if err.Error() != ErrMsgEmpty {
		t.Fatalf("caseEmpty: wrong err msg: expect=%s, but got=%s", ErrMsgEmpty, err)
	}
}

func TestIsChangedCaseOK(t *testing.T) {
	t.Parallel()

	var err error

	tempFile, err := NewTempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer Clean(tempFile)
	_, err = tempFile.IsChanged()
	if err.Error() != ErrMsgNotChanged {
		t.Fatalf("temporary file initialize must be non-changed status: expect=%s, but got=%s", ErrMsgNotChanged, err)
	}

	var caseOK = []struct {
		before string
		after  string
	}{
		{
			before: "foo",
			after:  "abc",
		},
		{
			before: "foo\n",
			after:  "abc\n",
		},
	}
	for i, test := range caseOK {
		tempFile.contents[previous], tempFile.contents[latest] =
			[]byte(test.before), []byte(test.after)
		changed, err := tempFile.IsChanged()
		if err != nil {
			t.Fatalf("caseOK: #%d: wrong err status: expect=nil, but got=%s", i+1, err)
		}
		if !changed {
			t.Fatalf("caseOK: #%d: wrong changed status: expect=%t, but got=%t", i+1, true, changed)
		}
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input  string
		expect string
	}{
		{
			input:  "test 1",
			expect: "test 1",
		},
		{
			input:  "test2\ntest2",
			expect: "test2\ntest2",
		},
		{
			input: `test3
test3
test3`,
			expect: "test3\ntest3\ntest3",
		},
	}

	for i, c := range testCases {
		tempFile, err := NewTempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer Clean(tempFile)
		tempFile.contents[latest] = []byte(c.input)
		output := tempFile.String()
		if output != c.expect {
			t.Fatalf("wrong test [%d]: expected=%s, but got=%s", i+1, c.expect, output)
		}
	}
}

func TestByte(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input  string
		expect []byte
	}{
		{
			input:  "test 1",
			expect: []byte("test 1"),
		},
		{
			input:  "test2\ntest2",
			expect: []byte("test2\ntest2"),
		},
		{
			input: `test3
test3
test3`,
			expect: []byte("test3\ntest3\ntest3"),
		},
	}

	for i, c := range testCases {
		tempFile, err := NewTempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer Clean(tempFile)
		tempFile.contents[latest] = []byte(c.input)
		output := tempFile.Byte()
		if !bytes.Equal(output, c.expect) {
			t.Fatalf("wrong test [%d]: expected=%s, but got=%s", i+1, c.expect, output)
		}
	}
}
