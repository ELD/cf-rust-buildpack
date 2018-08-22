package CommandExt

import (
	"io"
	"os/exec"
	"strings"
)

type CommandExt struct{}

func (ce *CommandExt) ExecuteWithPipe(dir string, stdout io.Writer, stderr io.Writer, command1, command2 string) error {
	commandArray1 := strings.Split(command1, " ")
	commandArray2 := strings.Split(command2, " ")

	c1 := exec.Command(commandArray1[0], commandArray1[1:]...)
	c2 := exec.Command(commandArray2[0], commandArray2[1:]...)

	c1.Dir = dir
	c2.Dir = dir

	pipe, _ := c1.StdoutPipe()
	defer pipe.Close()

	c2.Stdin = pipe
	c2.Stdout = stdout
	c2.Stderr = stderr

	c1.Start()

	err := c2.Run()

	return err
}
