// Copyright 2025 The MathWorks, Inc.

// This package intentionally does not use the LoggerFactory, or any logging framework.
// All messages are read and written as raw strings, so that the main MCP Server process can log them correctly.

package watchdog

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/stdio"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
)

type LoggerFactory interface {
	GetGlobalLogger() entities.Logger
}

type OSLayer interface {
	Getppid() int
	Stdin() io.Reader
	Stdout() io.Writer
	Stderr() io.Writer
}

type ProcessHandler interface {
	WatchProcessAndGetTerminationChan(processPid int) <-chan struct{}
	KillProcess(processPid int) error
}

type OSSignaler interface {
	InterruptSignalChan() <-chan os.Signal
}

type TransportFactory interface {
	NewReceiver(osStdio entities.OSStdio) (transport.Receiver, error)
}

type Watchdog struct {
	logger           entities.Logger
	osLayer          OSLayer
	processHandler   ProcessHandler
	osSignaler       OSSignaler
	transportFactory TransportFactory

	parentPID         int
	processPIDsToKill map[int]struct{}
	lock              *sync.Mutex
}

func New(
	loggerFactory LoggerFactory,
	osLayer OSLayer,
	processHandler ProcessHandler,
	osSignaler OSSignaler,
	transportFactory TransportFactory,
) *Watchdog {
	return &Watchdog{
		logger:           loggerFactory.GetGlobalLogger(),
		osLayer:          osLayer,
		processHandler:   processHandler,
		osSignaler:       osSignaler,
		transportFactory: transportFactory,

		processPIDsToKill: make(map[int]struct{}),
		lock:              new(sync.Mutex),
	}
}

func (w *Watchdog) StartAndWaitForCompletion(_ context.Context) error {
	receiver, err := w.transportFactory.NewReceiver(stdio.NewOSStdio(
		w.osLayer.Stdin(),
		w.osLayer.Stdout(),
		w.osLayer.Stderr(),
	))
	if err != nil {
		return err
	}

	w.logger.Info("Watchdog process has started")
	defer w.logger.Info("Watchdog process has exited")

	w.parentPID = w.osLayer.Getppid()

	shutdownMessageProcessingC := make(chan struct{})
	defer close(shutdownMessageProcessingC)

	shutdownC := make(chan struct{})
	go func() {
		c := receiver.C()
		for {
			select {
			case <-shutdownMessageProcessingC:
				return
			case rawMessage, ok := <-c:
				if !ok {
					w.logger.Error("Receiver channel closed unexpectedly")
					return
				}
				if abort := w.processIncomingMessage(rawMessage); abort {
					close(shutdownC)
					return
				}
			}
		}
	}()

	select {
	case <-shutdownC:
		defer func() {
			w.logger.Debug("Graceful shutdown completed")
			err := receiver.SendGracefulShutdownCompleted()
			if err != nil {
				w.logger.WithError(err).Error("Failed to send graceful shutdown completed signal")
			}
		}()
		w.logger.Debug("Graceful shutdown signal received")

	case <-w.processHandler.WatchProcessAndGetTerminationChan(w.parentPID):
		w.logger.Debug("Lost connection to parent, shutting down")

	case <-w.osSignaler.InterruptSignalChan():
		w.logger.Debug("Received unexpected graceful shutdown OS signal")
	}

	w.terminateAllProcesses()

	return nil
}

func (w *Watchdog) processIncomingMessage(rawMessage transport.Message) (abort bool) {
	w.lock.Lock()
	defer w.lock.Unlock()

	abort = false

	switch message := rawMessage.(type) {
	case transport.ProcessToKill:
		w.logger.
			With("pid", message.PID).
			Info("Adding process to kill")
		w.processPIDsToKill[message.PID] = struct{}{}
	case transport.Shutdown:
		abort = true
	}

	return
}

func (w *Watchdog) terminateAllProcesses() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.logger.
		With("count", len(w.processPIDsToKill)).
		Info("Trying to terminate children")

	for pid := range w.processPIDsToKill {
		w.logger.
			With("pid", pid).
			Debug("Killing process")

		if err := w.processHandler.KillProcess(pid); err != nil {
			w.logger.
				WithError(err).
				With("pid", pid).
				Error("Failed to kill child")
		}
	}
}
