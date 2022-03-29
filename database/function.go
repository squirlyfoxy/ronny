package database

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

type Function struct {
	Name      string
	StartLine int
	EndLine   int
}

func (f *Function) Exec() (string, error) {
	//Execute the function
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("./db/bin/" + f.Name + ".exe")
	} else {
		cmd = exec.Command("./db/bin/" + f.Name)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	//Get the result
	result := out.String()

	//Remove the last \n
	result = strings.TrimSuffix(result, "\n")

	return result, nil
}
