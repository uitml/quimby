package cli

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func Editor(in []byte) ([]byte, error) {
	// First test: open a temporary file with VI and read the saved file.
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	defer tmp.Close()

	_, err = tmp.Write(in)
	if err != nil {
		return nil, err
	}

	// Open the file in VI and read the result
	command := exec.Command("vi", tmp.Name())
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	err = command.Run()
	if err != nil {
		return nil, err
	}

	// Process the file and apply the values
	res, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		return nil, err
	}

	return res, nil
}
