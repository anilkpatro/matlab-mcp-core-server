// Copyright 2025 The MathWorks, Inc.

package logger_test

import (
	"path/filepath"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	loggermocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/logger"
	osfacademocks "github.com/matlab/matlab-mcp-core-server/mocks/facades/osfacade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFactory_HappyPath(t *testing.T) {
	for _, logLevel := range []entities.LogLevel{
		entities.LogLevelDebug,
		entities.LogLevelInfo,
		entities.LogLevelWarn,
		entities.LogLevelError,
	} {
		t.Run(string(logLevel), func(t *testing.T) {
			// Arrange
			mockConfig := &loggermocks.MockConfig{}
			defer mockConfig.AssertExpectations(t)

			mockDirectory := &loggermocks.MockDirectory{}
			defer mockDirectory.AssertExpectations(t)

			mockFilenameFactory := &loggermocks.MockFilenameFactory{}
			defer mockFilenameFactory.AssertExpectations(t)

			mockOSLayer := &loggermocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockConfig.EXPECT().
				LogLevel().
				Return(logLevel).
				Once()

			mockConfig.EXPECT().
				WatchdogMode().
				Return(false).
				Once()

			expectedBaseDir := "/some/directory"
			mockDirectory.EXPECT().
				BaseDir().
				Return(expectedBaseDir).
				Once()

			expectedSuffix := "123"
			mockDirectory.EXPECT().
				ID().
				Return(expectedSuffix).
				Once()

			expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
			mockFilenameFactory.EXPECT().
				FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
				Return(expectedLogFile).
				Once()

			mockLogFile := &osfacademocks.MockFile{}
			mockOSLayer.EXPECT().
				Create(expectedLogFile).
				Return(mockLogFile, nil).
				Once()

			// Act
			factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)

			// Assert
			require.NoError(t, err, "NewFactory should not return an error for valid log level")
			assert.NotNil(t, factory, "Factory should not be nil")
		})
	}
}

func TestNewFactory_LogFileCreateError(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("info").
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	expectedBaseDir := "/some/directory"
	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	expectedSuffix := "123"
	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	expectedError := assert.AnError
	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(nil, expectedError).
		Once()

	// Act
	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, factory, "Factory should not be nil")
}
func TestNewFactory_InvalidLogLevel(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("invalid").
		Once()

		// Act
	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)

	// Assert
	require.Error(t, err, "NewFactory should return an error for invalid log level")
	assert.Contains(t, err.Error(), "unknown log level", "Error message should indicate invalid log level")
	assert.Nil(t, factory, "Factory should be nil when error occurs")
}

func TestFactory_NewMCPSessionLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("info").
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	expectedBaseDir := "/some/directory"
	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	expectedSuffix := "123"
	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockLogFile := &osfacademocks.MockFile{}
	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)
	require.NoError(t, err, "Factory creation should not fail")

	// Act
	logger := factory.NewMCPSessionLogger(&mcp.ServerSession{})

	// Assert
	assert.NotNil(t, logger, "Logger should not be nil")
}

func TestFactory_GetGlobalLogger_HappyPath(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("debug").
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	expectedBaseDir := "/some/directory"
	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	expectedSuffix := "123"
	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockLogFile := &osfacademocks.MockFile{}
	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)
	require.NoError(t, err, "Factory creation should not fail")

	// Act
	logger := factory.GetGlobalLogger()

	// Assert
	assert.NotNil(t, logger, "Global logger should not be nil")
}

func TestFactory_GetGlobalLogger_UsesWatchdogLogFileInWatchdogMode(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("debug").
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(true).
		Once()

	expectedBaseDir := "/some/directory"
	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	expectedSuffix := "123"
	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	expectedLogFile := filepath.Join(expectedBaseDir, "watchdog.log")
	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.WatchdogLogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockLogFile := &osfacademocks.MockFile{}
	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)
	require.NoError(t, err, "Factory creation should not fail")

	// Act
	logger := factory.GetGlobalLogger()

	// Assert
	assert.NotNil(t, logger, "Global logger should not be nil")
}

func TestFactory_GetGlobalLogger_IsSingleton(t *testing.T) {
	// Arrange
	mockConfig := &loggermocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	mockDirectory := &loggermocks.MockDirectory{}
	defer mockDirectory.AssertExpectations(t)

	mockFilenameFactory := &loggermocks.MockFilenameFactory{}
	defer mockFilenameFactory.AssertExpectations(t)

	mockOSLayer := &loggermocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockConfig.EXPECT().
		LogLevel().
		Return("warn").
		Once()

	mockConfig.EXPECT().
		WatchdogMode().
		Return(false).
		Once()

	expectedBaseDir := "/some/directory"
	mockDirectory.EXPECT().
		BaseDir().
		Return(expectedBaseDir).
		Once()

	expectedSuffix := "123"
	mockDirectory.EXPECT().
		ID().
		Return(expectedSuffix).
		Once()

	expectedLogFile := filepath.Join(expectedBaseDir, "server.log")
	mockFilenameFactory.EXPECT().
		FilenameWithSuffix(filepath.Join(expectedBaseDir, logger.LogFileName), logger.LogFileExt, expectedSuffix).
		Return(expectedLogFile).
		Once()

	mockLogFile := &osfacademocks.MockFile{}
	mockOSLayer.EXPECT().
		Create(expectedLogFile).
		Return(mockLogFile, nil).
		Once()

	factory, err := logger.NewFactory(mockConfig, mockDirectory, mockFilenameFactory, mockOSLayer)
	require.NoError(t, err, "Factory creation should not fail")

	// Act
	logger1 := factory.GetGlobalLogger()
	logger2 := factory.GetGlobalLogger()

	// Assert
	assert.NotNil(t, logger1, "First global logger should not be nil")
	assert.NotNil(t, logger2, "Second global logger should not be nil")
	assert.Same(t, logger1, logger2, "Global logger should be a singleton")
}
