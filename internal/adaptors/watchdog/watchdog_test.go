// Copyright 2025 The MathWorks, Inc.

package watchdog_test

import (
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	watchdogmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/watchdog"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	transportmocks "github.com/matlab/matlab-mcp-core-server/mocks/watchdog/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Assert
	assert.NotNil(t, watchdogInstance, "Watchdog instance should not be nil")
}

func TestWatchdog_Start_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	require.NoError(t, err, "Start should not return an error")
}

func TestWatchdog_Start_TransportFactoryError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(nil, expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the transport factory error")
}

func TestWatchdog_Start_WatchdogProcessStartError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Act
	err := watchdogInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the watchdog process start error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	debugMessageC := make(chan string)
	defer close(debugMessageC)
	errorMessageC := make(chan string)
	defer close(errorMessageC)
	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Start the watchdog first
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)

	// Assert
	require.NoError(t, err, "RegisterProcessPIDWithWatchdog should not return an error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_WaitsIfNotStarted(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	debugMessageC := make(chan string)
	defer close(debugMessageC)
	errorMessageC := make(chan string)
	defer close(errorMessageC)
	expectedPID := 12345

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)
	}()

	select {
	case <-errC:
		t.Fatal("RegisterProcessPIDWithWatchdog should block until Start is called")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: RegisterProcessPIDWithWatchdog blocks
	}

	// Start after we've tried to register a PID
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	require.NoError(t, <-errC, "RegisterProcessPIDWithWatchdog should not return an error")
}

func TestWatchdog_RegisterProcessPIDWithWatchdog_SendProcessPIDError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	debugMessageC := make(chan string)
	defer close(debugMessageC)
	errorMessageC := make(chan string)
	defer close(errorMessageC)
	expectedPID := 12345
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendProcessPID(expectedPID).
		Return(expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Start the watchdog first
	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.RegisterProcessPIDWithWatchdog(expectedPID)

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the SendProcessPID error")
}

func TestWatchdog_Stop_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendStop().
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.Stop()

	// Assert
	assert.NoError(t, err, "Stop should not return an error")
}

func TestWatchdog_Stop_StopErrors(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendStop().
		Return(expectedError).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	// Act
	err = watchdogInstance.Stop()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the Stop error")
}

func TestWatchdog_Stop_WaitsIfNotStarted(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockWatchdogProcess := &watchdogmocks.MockWatchdogProcess{}
	defer mockWatchdogProcess.AssertExpectations(t)

	mockTransportFactory := &watchdogmocks.MockTransportFactory{}
	defer mockTransportFactory.AssertExpectations(t)

	mockLoggerFactory := &watchdogmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockTransportClient := &transportmocks.MockClient{}
	defer mockTransportClient.AssertExpectations(t)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockWatchdogProcess.EXPECT().
		Stdio().
		Return(mockSubProcessStdio).
		Once()

	mockTransportFactory.EXPECT().
		NewClient(mockSubProcessStdio).
		Return(mockTransportClient, nil).
		Once()

	mockWatchdogProcess.EXPECT().
		Start().
		Return(nil).
		Once()

	mockTransportClient.EXPECT().
		SendStop().
		Return(nil).
		Once()

	watchdogInstance := watchdog.New(
		mockWatchdogProcess,
		mockTransportFactory,
		mockLoggerFactory,
	)

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- watchdogInstance.Stop()
	}()

	select {
	case <-errC:
		t.Fatal("Stop should block until started")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: Stop blocks until started
	}

	err := watchdogInstance.Start()
	require.NoError(t, err, "Start should not return an error")

	assert.NoError(t, <-errC, "Stop should not return an error")
}
