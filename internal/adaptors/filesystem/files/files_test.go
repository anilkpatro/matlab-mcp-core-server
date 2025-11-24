// Copyright 2025 The MathWorks, Inc.

package files_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/filesystem/files"
	filesmock "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/filesystem/files"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	// Act
	factory := files.NewFactory(mockOSLayer)

	// Assert
	assert.NotNil(t, factory, "Factory instance should not be nil")
}

func TestFactory_CreateFileWithUniqueSuffix_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFile := &osfacademocks.MockFile{}
	defer mockFile.AssertExpectations(t)

	dir := filepath.Join("/tmp")
	pattern := "testfile"
	expectedFileName := filepath.Join("/tmp/testfile-12345")
	expectedSuffix := "12345"

	mockOSLayer.EXPECT().
		CreateTemp(dir, "testfile-*").
		Return(mockFile, nil).
		Once()

	mockFile.EXPECT().
		Name().
		Return(expectedFileName).
		Once()

	mockFile.EXPECT().
		Close().
		Return(nil).
		Once()

	factory := files.NewFactory(mockOSLayer)

	// Act
	fileName, suffix, err := factory.CreateFileWithUniqueSuffix(filepath.Join(dir, pattern), "")

	// Assert
	require.NoError(t, err, "CreateFileWithUniqueSuffix should not return an error")
	assert.Equal(t, expectedFileName, fileName, "File name should match expected")
	assert.Equal(t, expectedSuffix, suffix, "Suffix should match expected")
}

func TestFactory_CreateFileWithUniqueSuffix_PatternAlreadyHasSeparator(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFile := &osfacademocks.MockFile{}
	defer mockFile.AssertExpectations(t)

	dir := filepath.Join("/tmp")
	pattern := "testfile-"
	expectedFileName := filepath.Join("/tmp/testfile-67890")
	expectedSuffix := "67890"

	mockOSLayer.EXPECT().
		CreateTemp(dir, "testfile-*").
		Return(mockFile, nil).
		Once()

	mockFile.EXPECT().
		Name().
		Return(expectedFileName).
		Once()

	mockFile.EXPECT().
		Close().
		Return(nil).
		Once()

	factory := files.NewFactory(mockOSLayer)

	// Act
	fileName, suffix, err := factory.CreateFileWithUniqueSuffix(filepath.Join(dir, pattern), "")

	// Assert
	require.NoError(t, err, "CreateFileWithUniqueSuffix should not return an error")
	assert.Equal(t, expectedFileName, fileName, "File name should match expected")
	assert.Equal(t, expectedSuffix, suffix, "Suffix should match expected")
}

func TestFactory_CreateFileWithUniqueSuffix_FileWithExtension(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFile := &osfacademocks.MockFile{}
	defer mockFile.AssertExpectations(t)

	dir := filepath.Join("/var/log/app")
	pattern := "logfile"
	extension := ".log"
	expectedFileName := filepath.Join("/var/log/app/logfile-999.log")
	expectedSuffix := "999"

	mockOSLayer.EXPECT().
		CreateTemp(dir, "logfile-*.log").
		Return(mockFile, nil).
		Once()

	mockFile.EXPECT().
		Name().
		Return(expectedFileName).
		Once()

	mockFile.EXPECT().
		Close().
		Return(nil).
		Once()

	factory := files.NewFactory(mockOSLayer)

	// Act
	fileName, suffix, err := factory.CreateFileWithUniqueSuffix(filepath.Join(dir, pattern), extension)

	// Assert
	require.NoError(t, err, "CreateFileWithUniqueSuffix should not return an error")
	assert.Equal(t, expectedFileName, fileName, "File name should match expected")
	assert.Equal(t, expectedSuffix, suffix, "Suffix should match expected")
}

func TestFactory_CreateFileWithUniqueSuffix_CreateTempError(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	dir := filepath.Join("/tmp")
	pattern := "testfile"
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		CreateTemp(dir, "testfile-*").
		Return(nil, expectedError).
		Once()

	factory := files.NewFactory(mockOSLayer)

	// Act
	fileName, suffix, err := factory.CreateFileWithUniqueSuffix(filepath.Join(dir, pattern), "")

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the CreateTemp error")
	assert.Empty(t, fileName, "File name should be empty on error")
	assert.Empty(t, suffix, "Suffix should be empty on error")
}

func TestFactory_CreateFileWithUniqueSuffix_CloseError(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockFile := &osfacademocks.MockFile{}
	defer mockFile.AssertExpectations(t)

	dir := filepath.Join("/tmp")
	pattern := "testfile"
	expectedError := assert.AnError

	mockOSLayer.EXPECT().
		CreateTemp(dir, "testfile-*").
		Return(mockFile, nil).
		Once()

	mockFile.EXPECT().
		Close().
		Return(expectedError).
		Once()

	factory := files.NewFactory(mockOSLayer)

	// Act
	fileName, suffix, err := factory.CreateFileWithUniqueSuffix(filepath.Join(dir, pattern), "")

	// Assert
	require.ErrorIs(t, err, expectedError, "Error should be the Close error")
	assert.Empty(t, fileName, "File name should be empty on error")
	assert.Empty(t, suffix, "Suffix should be empty on error")
}

func TestFactory_FilenameWithSuffix_RegularFile(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	factory := files.NewFactory(mockOSLayer)

	fileName := "/path/to/file"
	extension := ".txt"
	suffix := "12345"
	expectedResult := "/path/to/file-12345.txt"

	// Act
	result := factory.FilenameWithSuffix(fileName, extension, suffix)

	// Assert
	assert.Equal(t, expectedResult, result, "Filename with suffix should match expected")
}

func TestFactory_FilenameWithSuffix_HiddenFile(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	factory := files.NewFactory(mockOSLayer)

	fileName := "/path/to/.hiddenfile"
	suffix := "67890"
	expectedResult := "/path/to/.hiddenfile-67890"

	// Act
	result := factory.FilenameWithSuffix(fileName, "", suffix)

	// Assert
	assert.Equal(t, expectedResult, result, "Hidden filename with suffix should match expected")
}

func TestFactory_FilenameWithSuffix_FileWithoutExtension(t *testing.T) {
	// Arrange
	mockOSLayer := &filesmock.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	factory := files.NewFactory(mockOSLayer)

	fileName := "/path/to/file"
	suffix := "abc"
	expectedResult := "/path/to/file-abc"

	// Act
	result := factory.FilenameWithSuffix(fileName, "", suffix)

	// Assert
	assert.Equal(t, expectedResult, result, "Filename without extension with suffix should match expected")
}
