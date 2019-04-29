package commander

import (
	"context"
	"os"
	"os/exec"
	"syscall"
)

type (
	Command interface {
		Run() error
	}

	cmd struct {
		program       string
		args          []string
		env           []string
		stdout        *os.File
		stdin         *os.File
		stderr        *os.File
		logger        logger
		ctx           context.Context
		detached      bool
		workDir       string
		signalHandler signalHandler
	}

	logger interface {
		Debugf(format string, args ...interface{})
		Debug(args ...interface{})
	}

	signalHandler interface {
		Process()
		Push(chan os.Signal)
	}

	Options struct {
		Program       string
		Args          []string
		Env           []string
		Detached      bool
		Logger        logger
		WorkDir       string
		SignalHandler signalHandler
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

func (c *cmd) Run() error {
	var gracefulStop = make(chan os.Signal)
	c.signalHandler.Push(gracefulStop)
	c.signalHandler.Process()

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

	err := command.Start()
	if err != nil {
		return err
	}
	go func(pid int) {
		c.logger.Debugf("Started go rutine to handle signal for process PPID: %d\n", pid)
		sig := <-gracefulStop
		c.logger.Debugf("Killing process %d\n", pid)
		syscall.Kill(-pid, syscall.SIGKILL)
		gracefulStop <- sig
	}(command.Process.Pid)

	err = command.Wait()
	return err
}
