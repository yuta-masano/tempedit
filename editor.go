package tempedit

import (
	"os"
	"runtime"
)

type runner interface {
	run(filePath string) error
}

// Editor is an Editor struct to open a temporary file.
type Editor struct {
	runner
}

func defaultName() string {
	val := os.Getenv("EDITOR")
	if val == "" {
		if runtime.GOOS == "windows" {
			val = "notepad.exe"
		} else {
			val = "vi"
		}
	}
	return val
}

// NewEditor creates a new Editor to open a temporary file.
// If name is omitted, $EDITOR or default editor (Win: notepad.exe, *nix: vi)
// is used.
func NewEditor(name string) *Editor {
	if name == "" {
		name = defaultName()
	}
	return &Editor{newApplication(name)}
}
