package commander

import (
	"os"
	"os/exec"
)

type (
	Command interface {
		Run() error
	}

	cmd struct {
		program string
		args    []string
		env     []string
		stdout  *os.File
		stdin   *os.File
		stderr  *os.File
	}
)

func New(program string, args []string, env []string) Command {
	return &cmd{
		stderr:  os.Stderr,
		stdout:  os.Stdout,
		stdin:   os.Stdin,
		program: program,
		args:    args,
		env:     env,
	}
}

func (c *cmd) Run() error {
	command := exec.Command(c.program, c.args...)
	command.Env = append(os.Environ(), c.env...)
	command.Stdout = c.stdout
	command.Stdin = c.stdin
	command.Stderr = c.stderr
	err := command.Start()
	if err != nil {
		return err
	}
	err = command.Wait()
	return err
}
