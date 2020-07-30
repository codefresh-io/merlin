package commander

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/codefresh-io/go/logger"
)

type (
	Command interface {
		Run(context.Context) error
		DryRun() string
	}

	cmd struct {
		program  string
		args     []string
		env      []string
		stdout   *os.File
		stdin    *os.File
		stderr   *os.File
		logger   logger.Logger
		ctx      context.Context
		detached bool
		workDir  string
	}

	signalHandler interface {
		Process()
		Push(chan os.Signal)
	}

	Options struct {
		Program  string
		Args     []string
		Env      []string
		Detached bool
		Logger   logger.Logger
		WorkDir  string
	}
)

func New(opt *Options) Command {
	ctx := context.Background()
	c := &cmd{
		stderr:        os.Stderr,
		stdout:        os.Stdout,
		stdin:         os.Stdin,
		program:       opt.Program,
		args:          opt.Args,
		env:           opt.Env,
		ctx:           ctx,
		detached:      opt.Detached,
		logger:        opt.Logger,
		workDir:       opt.WorkDir,
		signalHandler: opt.SignalHandler,
	}
	return c
}

func (c *cmd) Run(context context.Context) error {
	command := exec.CommandContext(c.ctx, c.program, c.args...)
	command.Env = append(os.Environ(), c.env...)
	if c.workDir != "" {
		command.Dir = c.workDir
	}
	if !c.detached {
		command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	} else {
		c.logger.Debug("Running in detached mode")
	}
	command.Stdout = c.stdout
	command.Stdin = c.stdin
	command.Stderr = c.stderr
	c.logger.Debug("Starting command", "workdir", c.workDir, "program", c.program, "arguments", c.args, "envs", c.env)
	err := command.Start()
	if err != nil {
		return err
	}
	go func(pid int) {
		c.logger.Debug("Started go rutine to handle signal for process PPID", "pid", pid)
		select {
		case <-context.Done():
			c.logger.Debug("Killing process", "pid", pid)
			syscall.Kill(-pid, syscall.SIGKILL)
		}
	}(command.Process.Pid)

	err = command.Wait()
	return err
}

func (c *cmd) DryRun() string {
	return fmt.Sprintf("workdir: %s --- envs: %v --- cmd: %s %s", c.workDir, c.env, c.program, c.args)
}
