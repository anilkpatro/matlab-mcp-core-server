// Copyright 2025 The MathWorks, Inc.

package transport

type Message interface {
	seal()
}

type ProcessToKill struct {
	PID int
}

func (p ProcessToKill) seal() {}

type Shutdown struct{}

func (p Shutdown) seal() {}

type Client interface {
	SendProcessPID(processPID int) error
	SendStop() error
}

type Receiver interface {
	SendGracefulShutdownCompleted() error

	C() <-chan Message
}
