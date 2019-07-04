package tempedit

import (
	"os"
	"reflect"
	"testing"
)

func TestDefaultName(t *testing.T) {
	var (
		input  string
		expect string
		output string
		err    error
	)

	input = "" // NOTE: This project supposes to be developed on *nix environment.
	expect = "vi"
	err = os.Setenv("EDITOR", input)
	if err != nil {
		t.Fatal(err)
	}
	output = defaultName()
	if output != expect {
		t.Fatalf("wrong test [1]: expected=%s, but got=%s", expect, output)
	}

	input = "test"
	expect = "test"
	err = os.Setenv("EDITOR", input)
	if err != nil {
		t.Fatal(err)
	}
	output = defaultName()
	if output != expect {
		t.Fatalf("wrong test [2]: expected=%s, but got=%s", expect, output)
	}
}

func TestNewEditor(t *testing.T) {
	_ = os.Unsetenv("EDITOR")

	var (
		expect *Editor
		output *Editor
	)

	expect = &Editor{&application{name: "vi"}} // NOTE: This project supposes to be developed on *nix environment.
	output = NewEditor("")
	if !reflect.DeepEqual(*expect, *output) {
		t.Fatalf("wrong test [1]: expect=%v, but got=%v", expect, output)
	}

	expect = &Editor{&application{name: "test"}}
	output = NewEditor("test")
	if !reflect.DeepEqual(*expect, *output) {
		t.Fatalf("wrong test [2]: expect=%v, but got=%v", expect, output)
	}
}
