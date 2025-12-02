// Copyright 2025 The MathWorks, Inc.

package watchdog_test

import (
	"os"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/stdio"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog"
	transportmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Assert
	assert.NotNil(t, watchdogInstance, "Watchdog instance should not be nil")
}

func TestWatchdog_StartAndWatch_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockReceiver := &transportmocks.MockReceiver{}
	defer mockReceiver.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	expectedParentPID := 1234

	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	expectedPIDToKill := 123654
	messageC := make(chan transport.Message)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockTransportFactory.EXPECT().
		NewReceiver(stdio.NewOSStdio(mockStdin, mockStdout, mockStderr)).
		Return(mockReceiver, nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockReceiver.EXPECT().
		C().
		Return(messageC).
		Once()

	mockReceiver.EXPECT().
		SendGracefulShutdownCompleted().
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPIDToKill).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	messageC <- transport.ProcessToKill{PID: expectedPIDToKill}

	messageC <- transport.Shutdown{}

	// Assert
	require.NoError(t, <-errC, "StartAndWatch should not return an error on graceful shutdown")
}

func TestWatchdog_StartAndWatch_MulitplePIDs(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockReceiver := &transportmocks.MockReceiver{}
	defer mockReceiver.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	expectedParentPID := 1234

	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	expectedPIDToKill := 123654
	expectedSecondPIDToKill := 6587987
	messageC := make(chan transport.Message)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockTransportFactory.EXPECT().
		NewReceiver(stdio.NewOSStdio(mockStdin, mockStdout, mockStderr)).
		Return(mockReceiver, nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockReceiver.EXPECT().
		C().
		Return(messageC).
		Once()

	mockReceiver.EXPECT().
		SendGracefulShutdownCompleted().
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPIDToKill).
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedSecondPIDToKill).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	messageC <- transport.ProcessToKill{PID: expectedPIDToKill}
	messageC <- transport.ProcessToKill{PID: expectedSecondPIDToKill}

	messageC <- transport.Shutdown{}

	// Assert
	require.NoError(t, <-errC, "StartAndWatch should not return an error on graceful shutdown")
}

func TestWatchdog_StartAndWatch_ParentProcessTermination(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockReceiver := &transportmocks.MockReceiver{}
	defer mockReceiver.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	expectedParentPID := 1234

	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	expectedPIDToKill := 123654
	messageC := make(chan transport.Message)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockTransportFactory.EXPECT().
		NewReceiver(stdio.NewOSStdio(mockStdin, mockStdout, mockStderr)).
		Return(mockReceiver, nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockReceiver.EXPECT().
		C().
		Return(messageC).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPIDToKill).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	messageC <- transport.ProcessToKill{PID: expectedPIDToKill}

	close(parentTerminationC)

	// Assert
	require.NoError(t, <-errC, "StartAndWatch should not return an error when parent is terminated")
}

func TestWatchdog_StartAndWatch_OSSignalInterrupt(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockReceiver := &transportmocks.MockReceiver{}
	defer mockReceiver.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	expectedParentPID := 1234

	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	expectedPIDToKill := 123654
	messageC := make(chan transport.Message)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockTransportFactory.EXPECT().
		NewReceiver(stdio.NewOSStdio(mockStdin, mockStdout, mockStderr)).
		Return(mockReceiver, nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockReceiver.EXPECT().
		C().
		Return(messageC).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPIDToKill).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	messageC <- transport.ProcessToKill{PID: expectedPIDToKill}

	interruptSignalC <- os.Interrupt

	// Assert
	require.NoError(t, <-errC, "StartAndWatch should not return an error when receiving SIGINT signal")
}

func TestWatchdog_StartAndWatch_KillProcessError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockProcessHandler := &mocks.MockProcessHandler{}
	defer mockProcessHandler.AssertExpectations(t)

	mockOSSignaler := &mocks.MockOSSignaler{}
	defer mockOSSignaler.AssertExpectations(t)

	mockTransportFactory := &mocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockReceiver := &transportmocks.MockReceiver{}
	defer mockReceiver.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	expectedParentPID := 1234

	parentTerminationC := make(chan struct{})
	interruptSignalC := make(chan os.Signal, 1)

	expectedPIDToKill := 123654
	expectedSecondPIDToKill := 64687
	messageC := make(chan transport.Message)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSLayer.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSLayer.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockTransportFactory.EXPECT().
		NewReceiver(stdio.NewOSStdio(mockStdin, mockStdout, mockStderr)).
		Return(mockReceiver, nil).
		Once()

	mockOSLayer.EXPECT().
		Getppid().
		Return(expectedParentPID).
		Once()

	mockReceiver.EXPECT().
		C().
		Return(messageC).
		Once()

	mockReceiver.EXPECT().
		SendGracefulShutdownCompleted().
		Return(nil).
		Once()

	mockProcessHandler.EXPECT().
		WatchProcessAndGetTerminationChan(expectedParentPID).
		Return(parentTerminationC).
		Once()

	mockOSSignaler.EXPECT().
		InterruptSignalChan().
		Return(interruptSignalC).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedPIDToKill).
		Return(expectedError).
		Once()

	mockProcessHandler.EXPECT().
		KillProcess(expectedSecondPIDToKill).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockLoggerFactory,
		mockOSLayer,
		mockProcessHandler,
		mockOSSignaler,
		mockTransportFactory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.StartAndWaitForCompletion(t.Context())
	}()

	messageC <- transport.ProcessToKill{PID: expectedPIDToKill}
	messageC <- transport.ProcessToKill{PID: expectedSecondPIDToKill}

	messageC <- transport.Shutdown{}

	// Assert
	require.NoError(t, <-errC, "StartAndWatch should not return an error even if failing to kill a child on shutdown")
}
