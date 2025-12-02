// Copyright 2025 The MathWorks, Inc.

package server_test

import (
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/server"
	toolsmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockFirstTool := &toolsmocks.MockTool{}
	defer mockFirstTool.AssertExpectations(t)

	mockSecondTool := &toolsmocks.MockTool{}
	defer mockSecondTool.AssertExpectations(t)

	mockServerConfig := &mocks.MockServerConfig{}
	defer mockServerConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockServerConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	expectedMCPServer := server.NewMCPSDKServer(mockServerConfig)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return([]tools.Tool{mockFirstTool, mockSecondTool}).
		Once()

	mockFirstTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(nil).
		Once()

	mockSecondTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(nil).
		Once()

	// Act
	server, err := server.New(expectedMCPServer, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, server, "Server should not be nil")
}

func TestNew_AddToServerReturnsError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockTool := &toolsmocks.MockTool{}
	defer mockTool.AssertExpectations(t)

	mockServerConfig := &mocks.MockServerConfig{}
	defer mockServerConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedError := assert.AnError

	mockServerConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	expectedMCPServer := server.NewMCPSDKServer(mockServerConfig)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return([]tools.Tool{mockTool}).
		Once()

	mockTool.EXPECT().
		AddToServer(expectedMCPServer).
		Return(expectedError).
		Once()

	// Act
	server, err := server.New(expectedMCPServer, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Assert
	require.Error(t, err, "New should return an error")
	assert.Equal(t, expectedError, err, "Error should match expected error")
	assert.Empty(t, server, "Server should be nil when error occurs")
}

func TestNew_HandlesNoTools(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockServerConfig := &mocks.MockServerConfig{}
	defer mockServerConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockServerConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	expectedMCPServer := server.NewMCPSDKServer(mockServerConfig)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return(nil).
		Once()

	// Act
	server, err := server.New(expectedMCPServer, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)

	// Assert
	require.NoError(t, err, "New should not return an error")
	assert.NotNil(t, server, "Server should not be nil")
}

func TestServer_Run_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLifecycleSignaler := &mocks.MockLifecycleSignaler{}
	defer mockLifecycleSignaler.AssertExpectations(t)

	mockConfigurator := &mocks.MockMCPServerConfigurator{}
	defer mockConfigurator.AssertExpectations(t)

	mockServerConfig := &mocks.MockServerConfig{}
	defer mockServerConfig.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	mockServerConfig.EXPECT().
		Version().
		Return("1.0.0").
		Once()

	expectedMCPServer := server.NewMCPSDKServer(mockServerConfig)

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockConfigurator.EXPECT().
		GetToolsToAdd().
		Return(nil).
		Once()

	capturedShutdownFuncC := make(chan func() error)
	mockLifecycleSignaler.EXPECT().
		AddShutdownFunction(mock.AnythingOfType("func() error")).
		Run(func(shutdownFcn func() error) {
			capturedShutdownFuncC <- shutdownFcn
		}).
		Return().
		Once()

	server, err := server.New(expectedMCPServer, mockLoggerFactory, mockLifecycleSignaler, mockConfigurator)
	require.NoError(t, err)

	// The MCP STDIO transport will hijack os.Stdout, which will cause issues with code coverage reporting.
	// To avoid this, we replace the transport with an in memory transport.
	_, serverTransport := mcp.NewInMemoryTransports()
	server.SetServerTransport(serverTransport)

	errC := make(chan error)
	go func() {
		errC <- server.Run()
	}()

	capturedShutdownFunc := <-capturedShutdownFuncC

	// Act
	err = capturedShutdownFunc()

	// Assert
	require.NoError(t, err, "Shutdown function should not return an error")
	serverErr := <-errC
	require.NoError(t, serverErr, "Server run should exit without error after shutdown")
}
