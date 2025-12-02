// Copyright 2025 The MathWorks, Inc.

package basetool_test

import (
	"context"
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/mcp/tools/basetool"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestUnstructuredInput struct {
	Query string `json:"query"`
}

func TestNewToolWithUnstructuredContentOutput_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	const (
		toolName        = "test-unstructured-tool"
		toolTitle       = "Test Unstructured Tool"
		toolDescription = "A test tool for unstructured content"
	)

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	// Act
	tool := basetool.NewToolWithUnstructuredContent(
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

	expectedInputSchema, err := jsonschema.For[TestUnstructuredInput](&jsonschema.ForOptions{})
	require.NoError(t, err, "Input schema generation should succeed")
	inputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "Input schema generation should succeed")
	require.Equal(t, expectedInputSchema, inputSchema, "Input schema should match expected")
}

func TestToolWithUnstructuredContentOutput_AddToServer_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockAdder := &mocks.MockToolAdder[TestUnstructuredInput, any]{}
	defer mockAdder.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedServer := mcp.NewServer(&mcp.Implementation{}, &mcp.ServerOptions{})

	const (
		toolName        = "test-unstructured-tool"
		toolTitle       = "Test Unstructured Tool"
		toolDescription = "A test tool for unstructured content"
	)

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{
			TextContent: []string{"test response"},
		}, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
		toolName,
		toolTitle,
		toolDescription,
		mockLoggerFactory,
		handler,
	)

	toolInputSchema, err := tool.GetInputSchema()
	require.NoError(t, err, "GetInputSchema should not return an error")

	mockAdder.EXPECT().AddTool(
		expectedServer,
		&mcp.Tool{
			Name:         toolName,
			Title:        toolTitle,
			Description:  toolDescription,
			InputSchema:  toolInputSchema,
			OutputSchema: nil,
		},
		mock.Anything,
	)

	tool.SetToolAdder(mockAdder)

	// Act
	err = tool.AddToServer(expectedServer)

	// Assert
	require.NoError(t, err, "AddToServer should not return an error")
}

func TestToolWithUnstructuredContentOutput_Handler_HappyPath(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent:  []string{"text response"},
		ImageContent: []tools.PNGImageData{[]byte("image1")},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	textContent, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "First content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[0], textContent.Text, "Text content should match")

	imageContent, ok := result.Content[1].(*mcp.ImageContent)
	require.True(t, ok, "Second content should be image content")
	assert.Equal(t, "image/png", imageContent.MIMEType, "Image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[0]), imageContent.Data, "Image data should match")
}

func TestToolWithUnstructuredContentOutput_Handler_TextContentOnly(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent: []string{"response 1", "response 2"},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	textContent1, ok := result.Content[0].(*mcp.TextContent)
	require.True(t, ok, "First content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[0], textContent1.Text, "First text content should match")

	textContent2, ok := result.Content[1].(*mcp.TextContent)
	require.True(t, ok, "Second content should be text content")
	assert.Equal(t, expectedRichContent.TextContent[1], textContent2.Text, "Second text content should match")
}

func TestToolWithUnstructuredContentOutput_Handler_ImageContentOnly(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		ImageContent: []tools.PNGImageData{
			[]byte("image1"),
			[]byte("image2"),
		},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	require.Len(t, result.Content, 2, "Should have 2 content items")

	imageContent1, ok := result.Content[0].(*mcp.ImageContent)
	require.True(t, ok, "First content should be image content")
	assert.Equal(t, "image/png", imageContent1.MIMEType, "First image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[0]), imageContent1.Data, "First image data should match")

	imageContent2, ok := result.Content[1].(*mcp.ImageContent)
	require.True(t, ok, "Second content should be image content")
	assert.Equal(t, "image/png", imageContent2.MIMEType, "Second image MIME type should be PNG")
	assert.Equal(t, []byte(expectedRichContent.ImageContent[1]), imageContent2.Data, "Second image data should match")
}

func TestToolWithUnstructuredContentOutput_Handler_NoContent(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent:  []string{},
		ImageContent: []tools.PNGImageData{},
	}

	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	assert.Nil(t, output, "Output should be nil for unstructured content")
	require.NotNil(t, result, "Result should not be nil")
	assert.Empty(t, result.Content, "Content should be empty")
}

func TestToolWithUnstructuredContentOutput_Handler_UnstructuredHandlerError(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedError := assert.AnError
	mockSessionLogger := testutils.NewInspectableLogger()

	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		return tools.RichContent{}, expectedError
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	require.ErrorIs(t, err, expectedError, "Handler should return the expected error")
	assert.Nil(t, result, "Result should be nil when error occurs")
	assert.Nil(t, output, "Output should be nil when error occurs")
}

func TestToolWithUnstructuredContentOutput_Handler_ContextPropagation(t *testing.T) {
	// Arrange
	mockLoggerFactory := &mocks.MockLoggerFactory{}
	defer mockLoggerFactory.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()
	expectedSession := &mcp.ServerSession{}
	expectedInput := TestUnstructuredInput{Query: "test query"}
	expectedRichContent := tools.RichContent{
		TextContent: []string{"success"},
	}
	mockSessionLogger := testutils.NewInspectableLogger()

	contextReceived := make(chan context.Context, 1) // Buffering to avoid deadlock
	handler := func(ctx context.Context, logger entities.Logger, input TestUnstructuredInput) (tools.RichContent, error) {
		contextReceived <- ctx
		return expectedRichContent, nil
	}

	mockLoggerFactory.EXPECT().
		GetGlobalLogger().
		Return(mockLogger).
		Once()

	mockLoggerFactory.EXPECT().
		NewMCPSessionLogger(expectedSession).
		Return(mockSessionLogger).
		Once()

	tool := basetool.NewToolWithUnstructuredContent(
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
	assert.Equal(t, t.Context(), <-contextReceived, "Context should be propagated to handler")
}
