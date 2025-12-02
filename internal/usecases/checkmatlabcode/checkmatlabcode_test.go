// Copyright 2025 The MathWorks, Inc.

package checkmatlabcode_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	checkmatlabcodemocks "github.com/matlab/matlab-mcp-core-server/mocks/usecases/checkmatlabcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_HappyPath(t *testing.T) {
	// Arrange
	mockGlobalMATLAB := &entitiesmocks.MockGlobalMATLAB{}
	defer mockGlobalMATLAB.AssertExpectations(t)

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	// Act
	usecase := checkmatlabcode.New(mockPathValidator)

	// Assert
	assert.NotNil(t, usecase, "Usecase should not be nil")
}

func TestUsecase_Execute_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	validatedPath := filepath.Join("validated", "path", "to", "script.m")
	const expectedCheckCodeOutput = "L 5 (C 1-10): Variable 'x' might be unused."

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return(validatedPath, nil).
		Once()

	expectedEvalRequest := entities.EvalRequest{
		Code: "checkcode('" + validatedPath + "')",
	}

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: expectedCheckCodeOutput}, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, []string{expectedCheckCodeOutput}, response.CheckCodeOutput, "CheckCode output should match expected value")
}

func TestUsecase_Execute_HappyPath_OutputWithWhitespaceAndEmptyLines(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	validatedPath := filepath.Join("validated", "path", "to", "script.m")
	const expectedCheckCodeOutput = "  Line 1: Warning  \n\n  \n\nLine 3: Error\n   \n"
	expectedCleanedOutput := []string{"Line 1: Warning", "Line 3: Error"}

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return(validatedPath, nil).
		Once()

	expectedEvalRequest := entities.EvalRequest{
		Code: "checkcode('" + validatedPath + "')",
	}

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: expectedCheckCodeOutput}, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, expectedCleanedOutput, response.CheckCodeOutput, "CheckCode output should match expected value")
}

func TestUsecase_Execute_HappyPath_EmptyOutput(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	validatedPath := filepath.Join("validated", "path", "to", "script.m")
	const expectedCheckCodeOutput = "No issues found by checkcode"

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return(validatedPath, nil).
		Once()

	expectedEvalRequest := entities.EvalRequest{
		Code: "checkcode('" + validatedPath + "')",
	}

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: ""}, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, []string{expectedCheckCodeOutput}, response.CheckCodeOutput, "CheckCode output should match expected value")
}

func TestUsecase_Execute_HappyPath_PathWithSingleQuotes(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	validatedPath := filepath.Join("path", "with'quote", "script.m")
	const expectedCheckCodeOutput = "L 5 (C 1-10): Variable 'x' might be unused."

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return(validatedPath, nil).
		Once()

	expectedEvalRequest := entities.EvalRequest{
		Code: "checkcode('" + strings.ReplaceAll(validatedPath, "'", "''") + "')",
	}

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{ConsoleOutput: expectedCheckCodeOutput}, nil).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.NoError(t, err, "Execute should not return an error")
	assert.Equal(t, []string{expectedCheckCodeOutput}, response.CheckCodeOutput, "CheckCode output should match expected value")
}

func TestUsecase_Execute_PathValidationError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	expectedError := assert.AnError

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return("", expectedError).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}

func TestUsecase_Execute_EvalWithCaptureError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockPathValidator := &checkmatlabcodemocks.MockPathValidator{}
	defer mockPathValidator.AssertExpectations(t)

	mockClient := &entitiesmocks.MockMATLABSessionClient{}
	defer mockClient.AssertExpectations(t)

	checkcodeRequest := checkmatlabcode.Args{
		ScriptPath: filepath.Join("path", "to", "script.m"),
	}

	ctx := t.Context()
	validatedPath := filepath.Join("validated", "path", "to", "script.m")
	expectedError := assert.AnError

	mockPathValidator.EXPECT().
		ValidateMATLABScript(checkcodeRequest.ScriptPath).
		Return(validatedPath, nil).
		Once()

	expectedEvalRequest := entities.EvalRequest{
		Code: "checkcode('" + validatedPath + "')",
	}

	mockClient.EXPECT().
		EvalWithCapture(ctx, mockLogger.AsMockArg(), expectedEvalRequest).
		Return(entities.EvalResponse{}, expectedError).
		Once()

	usecase := checkmatlabcode.New(mockPathValidator)

	// Act
	response, err := usecase.Execute(ctx, mockLogger, mockClient, checkcodeRequest)

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the original error")
	assert.Empty(t, response, "Response should be empty when there's an error")
}
