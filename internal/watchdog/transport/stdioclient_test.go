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

func TestClient_SendProcessPID_HappyPath(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	expectedPID := 12345
	expectedMessageBytes := []byte(fmt.Sprintf("%d\n", expectedPID))

	mockStdin.EXPECT().
		Write(expectedMessageBytes).
		Return(len(expectedMessageBytes), nil).
		Once()

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdoutIsClosed)
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(10 * time.Millisecond)

	// Act
	err = client.SendProcessPID(expectedPID)

	// Assert
	require.NoError(t, err, "SendProcessPID should not return an error")
}

func TestClient_SendProcessPID_WriteError(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	expectedError := assert.AnError

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	mockStdin.EXPECT().
		Write(mock.Anything).
		Return(0, expectedError).
		Once()

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdoutIsClosed)
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(10 * time.Millisecond)

	// Act
	err = client.SendProcessPID(12345)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestClient_SendStop_HappyPath(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	sendShutdownCompleteMessage := make(chan struct{})
	shutdownCompleteMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownCompletedSignal))

	mockStdout.EXPECT().
		Read(mock.Anything).
		RunAndReturn(func(p []byte) (int, error) {
			<-sendShutdownCompleteMessage
			copy(p, shutdownCompleteMessageBytes)
			return len(shutdownCompleteMessageBytes), nil
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownSignal))

	mockStdin.EXPECT().
		Write(expectedMessageBytes).
		Return(len(expectedMessageBytes), nil).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(100 * time.Millisecond)

	// Act & Assert
	errC := make(chan error)
	go func() {
		errC <- client.SendStop()
	}()

	select {
	case <-errC:
		t.Fatal("Stop should not complete immediately")
	case <-time.After(10 * time.Millisecond):
		// Expected behavior: Stop does not complete immediately
	}

	close(sendShutdownCompleteMessage)

	assert.NoError(t, <-errC, "Stop should not return an error")
}

func TestClient_SendStop_WriteError(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdoutIsClosed)
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	expectedError := assert.AnError
	expectedMessageBytes := []byte(fmt.Sprintf("%s\n", transport.GracefulShutdownSignal))

	mockStdin.EXPECT().
		Write(expectedMessageBytes).
		Return(0, expectedError).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(10 * time.Millisecond)

	// Act
	err = client.SendStop()

	// Assert
	assert.ErrorIs(t, err, expectedError, "Stop should return an error when write fails")
}

func TestClient_DebugMessagesC_ReceivesMessages(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	expectedMessages := []string{
		"debug message 1",
		"debug message 2",
	}

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	for _, expectedMessage := range expectedMessages {
		mockStdout.EXPECT().
			Read(mock.Anything).
			RunAndReturn(func(p []byte) (int, error) {
				messageBytes := []byte(fmt.Sprintf("%s\n", expectedMessage))
				copy(p, messageBytes)
				return len(messageBytes), nil
			}).
			Once()
	}

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdoutIsClosed)
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(10 * time.Millisecond)

	// Act & Assert
	for _, expectedMessage := range expectedMessages {
		select {
		case msg := <-client.DebugMessagesC():
			assert.Equal(t, expectedMessage, msg, "Should receive debug message")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received debug message within timeout")
		}
	}
}

func TestClient_ErrorMessagesC_ReceivesMessages(t *testing.T) {
	// Arrange
	mockSubProcessStdio := &entitiesmocks.MockSubProcessStdio{}
	defer mockSubProcessStdio.AssertExpectations(t)

	expectedMessages := []string{
		"error message 1",
		"error message 2",
	}

	mockStdin := &entitiesmocks.MockWriter{}
	defer mockStdin.AssertExpectations(t)

	mockStdout := &entitiesmocks.MockReader{}
	defer mockStdout.AssertExpectations(t)

	mockStderr := &entitiesmocks.MockReader{}
	defer mockStderr.AssertExpectations(t)

	mockSubProcessStdio.EXPECT().
		Stdin().
		Return(mockStdin).
		Once()

	mockSubProcessStdio.EXPECT().
		Stdout().
		Return(mockStdout).
		Once()

	mockSubProcessStdio.EXPECT().
		Stderr().
		Return(mockStderr).
		Once()

	for _, expectedMessage := range expectedMessages {
		mockStderr.EXPECT().
			Read(mock.Anything).
			RunAndReturn(func(p []byte) (int, error) {
				messageBytes := []byte(fmt.Sprintf("%s\n", expectedMessage))
				copy(p, messageBytes)
				return len(messageBytes), nil
			}).
			Once()
	}

	blockUntilStdoutIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStdoutIsClosed
	}()

	blockUntilStderrIsClosed := make(chan struct{})
	defer func() {
		<-blockUntilStderrIsClosed
	}()

	mockStdout.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStdoutIsClosed)
		}).
		Once()

	mockStderr.EXPECT().
		Read(mock.Anything).
		Return(0, io.EOF).
		Run(func(p []byte) {
			close(blockUntilStderrIsClosed)
		}).
		Once()

	client, err := transport.NewStdioClient(mockSubProcessStdio)
	require.NoError(t, err)

	client.SetShutdownTimeout(10 * time.Millisecond)

	// Act & Assert
	for _, expectedMessage := range expectedMessages {
		select {
		case msg := <-client.ErrorMessagesC():
			assert.Equal(t, expectedMessage, msg, "Should receive error message")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Should have received received message within timeout")
		}
	}
}
