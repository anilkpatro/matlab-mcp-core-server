// Copyright 2025 The MathWorks, Inc.

package orchestrator_test

import (
	"os"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	orchestratormocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/orchestrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	//Act
	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Assert
	assert.NotNil(t, orchestratorInstance, "Orchestrator instance should not be nil")
}

func TestOrchestrator_StartAndWaitForCompletion_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	serverStarted := make(chan struct{})
	stopServer := make(chan struct{})
	defer close(stopServer)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Twice()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	// Server should run indefinitely (simulate with a blocking channel)
	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			close(serverStarted)
			<-stopServer
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	<-serverStarted

	sendInterruptSignal(interruptC)

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")
}

func TestOrchestrator_StartAndWaitForCompletion_ServerError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Twice()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(expectedError).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the server error")
}

func TestOrchestrator_StartAndWaitForCompletion_InitializeMATLABErrorDoesNotTriggerShutdown(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	expectedError := assert.AnError
	closeServerRoutine := make(chan struct{})
	isShutdownCalled := make(chan struct{})

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Twice()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		RunAndReturn(func() error {
			<-closeServerRoutine
			return nil
		}).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, expectedError).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(getInterruptChannel()).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Run(func() {
			close(isShutdownCalled)
		}).
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(expectedError).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	// Assert
	select {
	case <-isShutdownCalled:
		t.Fatal("RequestShutdown should not be called when MATLAB initialization fails")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: no shutdown request should occur
	}

	close(closeServerRoutine)
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error")
}

func TestOrchestrator_StartAndWaitForCompletion_WaitForShutdownToCompleteError(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockConfig.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t)

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()
	expectedError := assert.AnError

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Twice()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(true).
		Once()

	mockGlobalMATLABManager.EXPECT().
		Client(ctx, mockLogger.AsMockArg()).
		Return(nil, nil).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(expectedError).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Act
	errC := make(chan error)
	go func() {
		errC <- orchestratorInstance.StartAndWaitForCompletion(ctx)
	}()

	// Assert
	require.NoError(t, <-errC, "StartAndWaitForCompletion should not return an error on signal interrupt")

	// This is mostly optional
	logs := mockLogger.WarnLogs()

	fields, found := logs["MATLAB MCP Core Server application shutdown failed"]
	require.True(t, found, "Expected a warning log about shutdown failure")

	errField, found := fields["error"]
	require.True(t, found, "Expected an error field in the warning log")

	err, ok := errField.(error)
	require.True(t, ok, "Error field should be of type error")
	require.ErrorIs(t, err, expectedError, "Logged error should match the shutdown error")
}

func TestOrchestrator_runMATLABMCPServerMain_MultipleSession_HappyPath(t *testing.T) {
	// Arrange
	mockLifecycleSignaler := &orchestratormocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfig := &orchestratormocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockServer := &orchestratormocks.MockServer{}
	defer mockServer.AssertExpectations(t)

	mockWatchdogClient := &orchestratormocks.MockWatchdogClient{}
	defer mockWatchdogClient.AssertExpectations(t)

	mockLoggerFactory := &orchestratormocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockSignalLayer := &orchestratormocks.MockOSSignaler{}
	defer mockSignalLayer.AssertExpectations(t)

	mockGlobalMATLABManager := &orchestratormocks.MockGlobalMATLAB{}
	defer mockGlobalMATLABManager.AssertExpectations(t) // Implicit assertion here, Initialize should not be called

	mockDirectory := &orchestratormocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	ctx := t.Context()
	interruptC := getInterruptChannel()

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockConfig.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockDirectory.EXPECT().
		RecordToLogger(mockLogger.AsMockArg()).
		Return().
		Once()

	mockWatchdogClient.EXPECT().
		Start().
		Return(nil).
		Once()

	mockServer.EXPECT().
		Run().
		Return(nil).
		Once()

	mockConfig.EXPECT().
		UseSingleMATLABSession().
		Return(false).
		Once()

	mockSignalLayer.EXPECT().
		InterruptSignalChan().
		Return(interruptC).
		Once()

	mockLifecycleSignaler.EXPECT().
		RequestShutdown().
		Return().
		Once()

	mockLifecycleSignaler.EXPECT().
		WaitForShutdownToComplete().
		Return(nil).
		Once()

	mockWatchdogClient.EXPECT().
		Stop().
		Return(nil).
		Once()

	orchestratorInstance := orchestrator.New(
		mockLifecycleSignaler,
		mockConfig,
		mockServer,
		mockWatchdogClient,
		mockLoggerFactory,
		mockSignalLayer,
		mockGlobalMATLABManager,
		mockDirectory,
	)

	// Act
	err := orchestratorInstance.StartAndWaitForCompletion(ctx)

	// Assert
	require.NoError(t, err)
}

func getInterruptChannel() chan os.Signal {
	return make(chan os.Signal, 1)
}

func sendInterruptSignal(interruptC chan os.Signal) {
	interruptC <- os.Interrupt
}
