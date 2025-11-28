// Copyright 2025 The MathWorks, Inc.

package process_test

import (
	"errors"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	processmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/watchdog/process"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(mockStdout, nil).
		Once()

	mockCmd.EXPECT().
		StderrPipe().
		Return(mockStderr, nil).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, processInstance, "Process instance should not be nil")
}

func TestNew_StdinPipeError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"
	expectedError := errors.New("stdin pipe error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(nil, expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the stdin pipe error")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestNew_StdoutPipeError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"
	expectedError := errors.New("stdout pipe error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(nil, expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the stdout pipe error")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestNew_StderrPipeError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"
	expectedError := errors.New("stderr pipe error")

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(mockStdout, nil).
		Once()

	mockCmd.EXPECT().
		StderrPipe().
		Return(nil, expectedError).
		Once()

	// Act
	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the stderr pipe error")
	assert.Nil(t, processInstance, "Process instance should be nil on error")
}

func TestProcess_Start_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"

	// Setup mocks for New
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(mockStdout, nil).
		Once()

	mockCmd.EXPECT().
		StderrPipe().
		Return(mockStderr, nil).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	mockCmd.EXPECT().
		Start().
		Return(nil).
		Once()

	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)
	require.NoError(t, err, "New should not return an error")

	// Act
	err = processInstance.Start()

	// Assert
	assert.NoError(t, err, "Start should not return an error")
}

func TestProcess_Start_Error(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"
	expectedError := errors.New("start process error")

	// Setup mocks for New
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(mockStdout, nil).
		Once()

	mockCmd.EXPECT().
		StderrPipe().
		Return(mockStderr, nil).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	mockCmd.EXPECT().
		Start().
		Return(expectedError).
		Once()

	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)
	require.NoError(t, err, "New should not return an error")

	// Act
	err = processInstance.Start()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Error should be the start process error")
}

func TestProcess_Stdio_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockOSLayer := &processmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockLoggerFactory := &processmocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockDirectory := &processmocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockCmd := &osfacademocks.MockCmd{}
	defer mockCmd.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	programPath := "/path/to/program"
	args := []string{programPath, "arg1", "arg2"}
	baseDir := "/tmp/base/dir"
	serverID := "server-id"

	// Setup mocks for New
	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockOSLayer.EXPECT().
		Args().
		Return(args).
		Once()

	mockDirectory.EXPECT().
		BaseDir().
		Return(baseDir).
		Once()

	mockDirectory.EXPECT().
		ID().
		Return(serverID).
		Once()

	mockOSLayer.EXPECT().
		Command(programPath, []string{
			"--" + flags.WatchdogMode,
			"--" + flags.BaseDir, baseDir,
			"--" + flags.ServerInstanceID, serverID,
		}).
		Return(mockCmd).
		Once()

	mockCmd.EXPECT().
		StdinPipe().
		Return(mockStdin, nil).
		Once()

	mockCmd.EXPECT().
		StdoutPipe().
		Return(mockStdout, nil).
		Once()

	mockCmd.EXPECT().
		StderrPipe().
		Return(mockStderr, nil).
		Once()

	mockCmd.EXPECT().
		SetSysProcAttr(mock.Anything). // OS specific, not testable
		Once()

	processInstance, err := process.New(mockOSLayer, mockLoggerFactory, mockDirectory)
	require.NoError(t, err, "New should not return an error")

	// Act
	stdio := processInstance.Stdio()

	// Assert
	assert.Equal(t, mockStdin, stdio.Stdin(), "Stdin should match the mocked stdin pipe")
	assert.Equal(t, mockStdout, stdio.Stdout(), "Stdout should match the mocked stdout pipe")
	assert.Equal(t, mockStderr, stdio.Stderr(), "Stderr should match the mocked stderr pipe")
}
