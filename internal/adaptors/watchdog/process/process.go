// Copyright 2025 The MathWorks, Inc.

package process

import (
	"io"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/stdio"
)

type OSLayer interface {
	Args() []string
	Command(name string, arg ...string) osfacade.Cmd
}

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type Directory interface {
	BaseDir() string
	ID() string
}

type Process struct {
	osLayer OSLayer
	cmd     osfacade.Cmd
	logger  entities.Logger

	stdin  io.Writer
	stdout io.Reader
	stderr io.Reader
}

func New(
	osLayer OSLayer,
	loggerFactory LoggerFactory,
	directory Directory,
) (*Process, error) {
	logger := loggerFactory.GetGlobalLogger()

	programPath := osLayer.Args()[0]
	cmd := osLayer.Command(programPath,
		"--"+flags.WatchdogMode,
		"--"+flags.BaseDir, directory.BaseDir(),
		"--"+flags.ServerInstanceID, directory.ID())

	watchdogProcessStdin, err := cmd.StdinPipe()
	if err != nil {
		logger.WithError(err).Error("Failed to get stdin pipe for watchdog")
		return nil, err
	}

	watchodProcessStdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.WithError(err).Error("Failed to get stdout pipe for watchdog")
		return nil, err
	}

	watchdogProcessStderr, err := cmd.StderrPipe()
	if err != nil {
		logger.WithError(err).Error("Failed to get stderr pipe for watchdog")
		return nil, err
	}

	cmd.SetSysProcAttr(getSysProcAttrForDetachingAProcess())

	process := &Process{
		osLayer: osLayer,
		cmd:     cmd,
		logger:  logger,

		stdin:  watchdogProcessStdin,
		stdout: watchodProcessStdout,
		stderr: watchdogProcessStderr,
	}

	return process, nil
}

func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		p.logger.WithError(err).Error("Failed to start watchdog process")
		return err
	}

	return nil
}

func (p *Process) Stdio() entities.SubProcessStdio {
	return stdio.NewSubProcessStdio(p.stdin, p.stdout, p.stderr)
}
