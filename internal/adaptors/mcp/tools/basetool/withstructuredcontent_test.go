// Copyright 2025 The MathWorks, Inc.

package basetool_test

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestInput struct {
	Message string `json:"message"`
}

type TestOutput struct {
	Result string `json:"result"`
}

func TestNewToolWithStructuredContent_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	const (
		toolName        = "test-tool"
		toolTitle       = "Test Tool"
		toolDescription = "A test tool for unit testing"
	)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	tool := basetool.NewToolWithStructuredContent(
		toolName,
		toolTitle,
		toolDescription,
		mockLoggerFactory,
		handler,
	)

	// Assert
	assert.Equal(t, toolName, tool.Name(), "Tool name should match")
	assert.Equal(t, toolTitle, tool.Title(), "Tool title should match")
	assert.Equal(t, toolDescription, tool.Description(), "Tool description should match")

	expectedInputSchema, err := jsonschema.For[TestInput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Input schema generation should succeed")
	inputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "Input schema generation should succeed")
	require.Equal(t, expectedInputSchema, inputSchema, "Input schema should not be nil")

	expectedOutputSchema, err := jsonschema.For[TestOutput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Output schema generation should succeed")
	outputSchema, err := tool.GetOutputSchema()
	require.NoError(t, err, "Output schema generation should succeed")
	require.Equal(t, expectedOutputSchema, outputSchema, "Output schema should not be nil")
}

func TestToolWithStructuredContentOutput_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockAdder := &mocks.MockToolAdder[TestInput, TestOutput]{}
	defer mockAdder.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	const (
		toolName        = "test-tool"
		toolTitle       = "Test Tool"
		toolDescription = "A test tool for unit testing"
	)

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "success"}, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		toolName,
		toolTitle,
		toolDescription,
		mockLoggerFactory,
		handler,
	)

	toolInputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "GetInputSchema should not return an error")

	toolOutputSchema, err := tool.GetOutputSchema()
	require.NoError(t, err, "GetOutputSchema should not return an error")

	expectedServer := mcp.NewServer(&mcp.Implementation{}, &mcp.ServerOptions{})

	mockAdder.EXPECT().AddTool(
		expectedServer,
		&mcp.Tool{
			Name:         toolName,
			Title:        toolTitle,
			Description:  toolDescription,
			InputSchema:  toolInputSchema,
			OutputSchema: toolOutputSchema,
		},
		mock.Anything,
	)

	tool.SetToolAdder(mockAdder)

	// Act
	err = tool.AddToServer(expectedServer)

	// Assert
	require.NoError(t, err, "AddToServer should not return an error")
}

func TestToolWithStructuredContentOutput_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedOutput := TestOutput{Result: "processed: test message"}
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{Result: "processed: " + input.Message}, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockGlobalLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Nil(t, result, "Result should be nil for structured content output")
	assert.Equal(t, expectedOutput, output, "Output should match expected output")
}

func TestToolWithStructuredContentOutput_Handler_StructuredHandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedError := assert.AnError
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		return TestOutput{}, expectedError
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockGlobalLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	result, output, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.ErrorIs(t, err, expectedError, "Handler should return an error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Empty(t, output, "Output should be zero value when error occurs")
}

func TestToolWithStructuredContentOutput_Handler_ContextPropagation(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockGlobalLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestInput{Message: "test message"}
	expectedOutput := TestOutput{Result: "success"}
	mockSessionLogger := testutils.NewInspectableLogger()
	var capturedContext context.Context

	handler := func(ctx context.Context, logger entities.Logger, input TestInput) (TestOutput, error) {
		capturedContext = ctx
		return expectedOutput, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockGlobalLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithStructuredContent(
		"test-tool",
		"Test Tool",
		"A test tool",
		mockLoggerFactory,
		handler,
	)

	req := &mcp.CallToolRequest{
		Session: expectedSession,
	}

	// Act
	_, _, err := tool.Handler()(t.Context(), req, expectedInput)

	// Assert
	require.NoError(t, err, "Handler should not return an error")
	assert.Equal(t, t.Context(), capturedContext, "Context should be propagated to handler")
}
