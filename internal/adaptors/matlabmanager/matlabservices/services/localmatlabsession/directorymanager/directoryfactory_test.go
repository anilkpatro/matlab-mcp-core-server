// Copyright 2025 The MathWorks, Inc.

package directorymanager_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectory := &mocks.MockApplicationDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	// Act
	factory := directorymanager.NewFactory(mockOSLayer, mockApplicationDirectory, mockMATLABFiles)

	// Assert
	assert.NotNil(t, factory)
}

func TestDirectoryFactory_Create_HappyPath(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectory := &mocks.MockApplicationDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedSessionDir := filepath.Join("tmp", "matlab-session-12345")
	packageDir := filepath.Join(expectedSessionDir, "+matlab_mcp")
	expectedCertificateFile := filepath.Join(expectedSessionDir, "cert.pem")
	expectedCertificateKeyFile := filepath.Join(expectedSessionDir, "cert.key")
	expectedMATLABFiles := map[string][]byte{
		"initializeMCP.m": []byte("some content"),
		"eval.m":          []byte("some other content"),
	}

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(expectedSessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	for fileName, fileContent := range expectedMATLABFiles {
		filePath := filepath.Join(packageDir, fileName)
		mockOSLayer.EXPECT().
			WriteFile(filePath, fileContent, os.FileMode(0o600)).
			Return(nil).
			Once()
	}

	factory := directorymanager.NewFactory(mockOSLayer, mockApplicationDirectory, mockMATLABFiles)

	// Act
	directory, err := factory.Create(mockLogger)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, directory)
	assert.Equal(t, expectedSessionDir, directory.Path())
	assert.Equal(t, expectedCertificateFile, directory.CertificateFile())
	assert.Equal(t, expectedCertificateKeyFile, directory.CertificateKeyFile())
}

func TestDirectoryFactory_Create_MkdirTempError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectory := &mocks.MockApplicationDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	expectedError := assert.AnError

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return("", expectedError).
		Once()

	factory := directorymanager.NewFactory(mockOSLayer, mockApplicationDirectory, mockMATLABFiles)

	// Act
	directory, err := factory.Create(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, directory)
}

func TestDirectoryFactory_Create_PackageDirectoryMkdirError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectory := &mocks.MockApplicationDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := "/tmp/matlab-session-12345"
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(expectedError).
		Once()

	factory := directorymanager.NewFactory(mockOSLayer, mockApplicationDirectory, mockMATLABFiles)

	// Act
	directory, err := factory.Create(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, directory)
}

func TestDirectoryFactory_Create_WriteFileError(t *testing.T) {
	// Arrange
	mockOSLayer := &mocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockApplicationDirectory := &mocks.MockApplicationDirectory{}
	defer mockApplicationDirectory.AssertExpectations(t)

	mockMATLABFiles := &mocks.MockMATLABFiles{}
	defer mockMATLABFiles.AssertExpectations(t)

	mockLogger := testutils.NewInspectableLogger()

	sessionDir := "/tmp/matlab-session-12345"
	packageDir := filepath.Join(sessionDir, "+matlab_mcp")
	expectedError := assert.AnError

	mockApplicationDirectory.EXPECT().
		CreateSubDir(mock.AnythingOfType("string")).
		Return(sessionDir, nil).
		Once()

	mockOSLayer.EXPECT().
		Mkdir(packageDir, os.FileMode(0o700)).
		Return(nil).
		Once()

	expectedFailingFileName := "initializeMCP.m"

	expectedMATLABFiles := map[string][]byte{
		expectedFailingFileName: []byte("some other content"),
	}

	mockMATLABFiles.EXPECT().
		GetAll().
		Return(expectedMATLABFiles).
		Once()

	mockOSLayer.EXPECT().
		WriteFile(filepath.Join(packageDir, expectedFailingFileName), expectedMATLABFiles[expectedFailingFileName], os.FileMode(0o600)).
		Return(expectedError).
		Once()

	factory := directorymanager.NewFactory(mockOSLayer, mockApplicationDirectory, mockMATLABFiles)

	// Act
	directory, err := factory.Create(mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, directory)
}
