// Copyright 2025 The MathWorks, Inc.

package transport_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReceiver_SendDebugMessage_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessage := "debug message"
	expectedMessageBytes := fmt.Appendf([]byte{}, "%s\n", expectedMessage)

	mockStdout.EXPECT().
		Write(expectedMessageBytes).
		Return(len(expectedMessageBytes), nil).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	receiver.SendDebugMessage(expectedMessage)
}

func TestReceiver_SendDebugMessage_StdoutWriteError(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessage := "debug message"
	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", expectedMessage))

	expectedError := assert.AnError

	mockStdout.EXPECT().
		Write(expectedMessageBytes).
		Return(0, expectedError).
		Once()

	expectedErrorMessage := expectedError.Error()
	expectedErrorMessageBytes := []byte(fmt.Sprintf("%s\n", expectedErrorMessage))

	mockStderr.EXPECT().
		Write(expectedErrorMessageBytes).
		Return(len(expectedErrorMessageBytes), nil).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	receiver.SendDebugMessage(expectedMessage)

	// Assert (all mock assertions)
}

func TestReceiver_SendErrorMessage_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessage := "error message"
	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", expectedMessage))

	mockStderr.EXPECT().
		Write(expectedMessageBytes).
		Return(len(expectedMessageBytes), nil).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	receiver.SendErrorMessage(expectedMessage)
}

func TestReceiver_SendErrorMessage_StderrWriteError(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessage := "error message"
	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", expectedMessage))

	expectedError := assert.AnError

	mockStderr.EXPECT().
		Write(expectedMessageBytes).
		Return(0, expectedError).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	receiver.SendErrorMessage(expectedMessage)

	// Assert - (Mock assertions and stderr write errors are silent)
}

func TestReceiver_C_ProcessToKill_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	expectedPIDs := []int{12345, 67890}

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	for _, expectedPID := range expectedPIDs {
		expectedMessageBytes := []byte(fmt.Sprintf("%d\n", expectedPID))

		mockStdin.EXPECT().
			Read(mock.Anything).
			RunAndReturn(func(p []byte) (int, error) {
				copy(p, expectedMessageBytes)
				return len(expectedMessageBytes), nil
			}).
			Once()
	}

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act & Assert
	for _, expectedPID := range expectedPIDs {
		select {
		case message := <-receiver.C():
			processToKill, ok := message.(transport.ProcessToKill)
			require.True(t, ok)
			assert.Equal(t, expectedPID, processToKill.PID, "Should receive expected PID")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received PID within timeout")
		}
	}
}

func TestReceiver_C_ProcessToKill_InvalidPID(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	invalidPID := "not_a_number"
	expectedValidPID := 123245

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", invalidPID))

	mockStdin.EXPECT().
		Read(mock.Anything).
		RunAndReturn(func(p []byte) (int, error) {
			copy(p, expectedMessageBytes)
			return len(expectedMessageBytes), nil
		}).
		Once()

	expectedValidMessageBytes := []byte(fmt.Sprintf("%d\n", expectedValidPID))

	mockStdin.EXPECT().
		Read(mock.Anything).
		RunAndReturn(func(p []byte) (int, error) {
			copy(p, expectedValidMessageBytes)
			return len(expectedValidMessageBytes), nil
		}).
		Once()

	mockStderr.EXPECT().
		Write(mock.Anything). // Don't care what error message it was
		Return(0, nil).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act & Assert
	select {
	case message := <-receiver.C():
		processToKill, ok := message.(transport.ProcessToKill)
		require.True(t, ok)
		assert.Equal(t, expectedValidPID, processToKill.PID, "Should receive expected PID")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Should have received PID within timeout")
	}
}

func TestReceiver_C_Shutdown_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownSignal))

	mockStdin.EXPECT().
		Read(mock.Anything).
		RunAndReturn(func(p []byte) (int, error) {
			copy(p, expectedMessageBytes)
			return len(expectedMessageBytes), nil
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act & Assert
	select {
	case message := <-receiver.C():
		assert.IsType(t, transport.Shutdown{}, message, "Should receive shutdown signal")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Should have received shutdown signal within timeout")
	}
}

func TestReceiver_SendGracefulShutdownCompleted_HappyPath(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownCompletedSignal))

	mockStdout.EXPECT().
		Write(expectedMessageBytes).
		Return(0, nil).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	err = receiver.SendGracefulShutdownCompleted()

	// Assert
	require.NoError(t, err, "GracefulShutdownCompleted should not return an error")
}

func TestReceiver_SendGracefulShutdownCompleted_WriteError(t *testing.T) {
	// Arrange
	mockOSStdio := &entitiesmocks.MockOSStdio{}
	defer mockOSStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockReader{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockWriter{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockWriter{}
	defer mockStderr.AssertExpectations(t)

	mockOSStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockOSStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockOSStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownCompletedSignal))
	expectedError := assert.AnError

	mockStdout.EXPECT().
		Write(expectedMessageBytes).
		Return(0, expectedError).
		Once()

	blockUntilStdinIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdinIsClosed
	}()

	mockStdin.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdinIsClosed)
		}).
		Once()

	receiver, err := transport.NewStdioReceiver(mockOSStdio)
	require.NoError(t, err)

	// Act
	err = receiver.SendGracefulShutdownCompleted()

	// Assert
	assert.ErrorIs(t, err, expectedError, "GracefulShutdownCompleted should return the write error")
}
