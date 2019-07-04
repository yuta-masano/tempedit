package tempedit

import (
	"os"
	"os/exec"
	"runtime"
)

type application struct {
	name string
}

func newApplication(name string) runner {
	return &application{name: name}
}

func (a *application) run(filePath string) error {
	editCmd := exec.Command(a.name, filePath)
	var stdin *os.File
	var err error
	if runtime.GOOS == "windows" {
		stdin, err = os.Open("CONIN$")
		if err != nil {
			panic(err)
		}
	} else {
		stdin = os.Stdin
	}
	editCmd.Stdin, editCmd.Stdout, editCmd.Stderr =
		stdin, os.Stdout, os.Stderr
	return editCmd.Run()
}
