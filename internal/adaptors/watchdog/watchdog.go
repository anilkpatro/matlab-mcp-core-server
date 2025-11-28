// Copyright 2025 The MathWorks, Inc.

package watchdog

import (
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
)

type WatchdogProcess interface {
	Start() error
	Stdio() entities.SubProcessStdio
}

type TransportFactory interface {
	NewClient(subProcessStdio entities.SubProcessStdio) (transport.Client, error)
}

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type Watchdog struct {
	logger entities.Logger

	transportFactory TransportFactory
	watchdogProcess  WatchdogProcess

	client transport.Client

	startedC chan struct{}
}

func New(
	watchdogProcess WatchdogProcess,
	transportFactory TransportFactory,
	loggerFactory LoggerFactory,
) *Watchdog {
	return &Watchdog{
		logger: loggerFactory.GetGlobalLogger(),

		transportFactory: transportFactory,
		watchdogProcess:  watchdogProcess,

		startedC: make(chan struct{}),
	}
}

func (w *Watchdog) Start() error {
	w.logger.Debug("Starting watchdog")

	client, err := w.transportFactory.NewClient(w.watchdogProcess.Stdio())
	if err != nil {
		return err
	}

	w.client = client

	err = w.watchdogProcess.Start()
	if err != nil {
		w.logger.WithError(err).Error("Failed to start watchdog process")
		return err
	}

	close(w.startedC)

	w.logger.Debug("Started watchdog")

	return nil
}

func (w *Watchdog) RegisterProcessPIDWithWatchdog(processPID int) error {
	<-w.startedC

	w.logger.With("pid", processPID).Debug("Adding child process to watchdog")
	return w.client.SendProcessPID(processPID)
}

func (w *Watchdog) Stop() error {
	<-w.startedC

	w.logger.Debug("Sending graceful shutdown signal to watchdog")
	return w.client.SendStop()
}
