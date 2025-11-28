// Copyright 2025 The MathWorks, Inc.

package config_test

import (
	"runtime/debug"
	"testing"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedConfig struct {
	versionMode                      bool
	disableTelemetry                 bool
	useSingleMATLABSession           bool
	logLevel                         entities.LogLevel
	preferredLocalMATLABRoot         string
	preferredMATLABStartingDirectory string
	baseDirectory                    string
	watchdogMode                     bool
	serverInstanceID                 string
}

func TestNew_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected expectedConfig
	}{
		{
			name: "default values",
			args: []string{},
			expected: expectedConfig{
				versionMode:                      false,
				disableTelemetry:                 false,
				useSingleMATLABSession:           true,
				logLevel:                         entities.LogLevelInfo,
				preferredLocalMATLABRoot:         "",
				preferredMATLABStartingDirectory: "",
				baseDirectory:                    "",
				watchdogMode:                     false,
				serverInstanceID:                 "",
			},
		},
		{
			name: "custom values",
			args: []string{
				"--version=true",
				"--disable-telemetry",
				"--use-single-matlab-session=false",
				"--log-level=debug",
				"--matlab-root=/root",
				"--initial-working-folder=/pref",
				"--log-folder=/logs",
				"--watchdog=true",
				"--server-instance-id=69878763",
			},
			expected: expectedConfig{
				versionMode:                      true,
				disableTelemetry:                 true,
				useSingleMATLABSession:           false,
				logLevel:                         entities.LogLevelDebug,
				preferredLocalMATLABRoot:         "/root",
				preferredMATLABStartingDirectory: "/pref",
				baseDirectory:                    "/logs",
				watchdogMode:                     true,
				serverInstanceID:                 "69878763",
			},
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			// Act
			cfg, err := config.New(mockOSLayer)

			// Assert
			require.NoError(t, err)
			require.NotNil(t, cfg, "Config should not be nil")

			assert.Equal(t, testConfig.expected.versionMode, cfg.VersionMode())
			assert.Equal(t, testConfig.expected.disableTelemetry, cfg.DisableTelemetry())
			assert.Equal(t, testConfig.expected.useSingleMATLABSession, cfg.UseSingleMATLABSession())
			assert.Equal(t, testConfig.expected.logLevel, cfg.LogLevel())
			assert.Equal(t, testConfig.expected.preferredLocalMATLABRoot, cfg.PreferredLocalMATLABRoot())
			assert.Equal(t, testConfig.expected.preferredMATLABStartingDirectory, cfg.PreferredMATLABStartingDirectory())
			assert.Equal(t, testConfig.expected.baseDirectory, cfg.BaseDir())
			assert.Equal(t, testConfig.expected.watchdogMode, cfg.WatchdogMode())
			assert.Equal(t, testConfig.expected.serverInstanceID, cfg.ServerInstanceID())
		})
	}
}

func TestConfig_Version_VersionSetByLDFLAGS(t *testing.T) {
	// Arrange
	expectedVersion := "v1.2.3"
	config.SetVersionLikeLDFLAGSWould(t, expectedVersion)

	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Args().
		Return([]string{"testprocess"}).
		Once()

	expectedPath := "/some/path/to/module"
	versionToNotBeUsed := "v0.0.0"
	mockOSLayer.EXPECT().
		ReadBuildInfo().
		Return(&debug.BuildInfo{
			Main: debug.Module{
				Path:    expectedPath,
				Version: versionToNotBeUsed,
			},
		}, true).
		Once()

	cfg, err := config.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	version := cfg.Version()

	// Assert
	require.Equal(t, expectedPath+" "+expectedVersion, version)
}

func TestConfig_Version_VersionUnsetByLDFLAGS(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Args().
		Return([]string{"testprocess"}).
		Once()

	expectedPath := "/some/path/to/module"
	expectedVersion := "v1.2.3"
	mockOSLayer.EXPECT().
		ReadBuildInfo().
		Return(&debug.BuildInfo{
			Main: debug.Module{
				Path:    expectedPath,
				Version: expectedVersion,
			},
		}, true).
		Once()

	cfg, err := config.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	version := cfg.Version()

	// Assert
	require.Equal(t, expectedPath+" "+expectedVersion, version)
}

func TestConfig_Version_VersionUnsetByLDFLAGS_ButReadBuildInfoNotOK(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Args().
		Return([]string{"testprocess"}).
		Once()

	expectedPath := "/some/path/to/module"
	expectedVersion := "v1.2.3"
	mockOSLayer.EXPECT().
		ReadBuildInfo().
		Return(&debug.BuildInfo{
			Main: debug.Module{
				Path:    expectedPath,
				Version: expectedVersion,
			},
		}, false).
		Once()

	cfg, err := config.New(mockOSLayer)
	require.NoError(t, err)

	// Act
	version := cfg.Version()

	// Assert
	require.Equal(t, expectedPath+" (devel)", version)
}

func TestConfig_DisableTelemetry_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: false,
		},
		{
			name:     "implicitly true",
			args:     []string{"--disable-telemetry"},
			expected: true,
		},
		{
			name:     "explicitly true",
			args:     []string{"--disable-telemetry=true"},
			expected: true,
		},
		{
			name:     "explicitly false",
			args:     []string{"--disable-telemetry=false"},
			expected: false,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.DisableTelemetry()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_UseSingleMATLABSession_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: true,
		},
		{
			name:     "explicitly true",
			args:     []string{"--use-single-matlab-session=true"},
			expected: true,
		},
		{
			name:     "explicitly false",
			args:     []string{"--use-single-matlab-session=false"},
			expected: false,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.UseSingleMATLABSession()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_PreferredLocalMATLABRoot_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Windows MATLAB path",
			args:     []string{"--matlab-root=C:\\Program Files\\MATLAB\\R2024b"},
			expected: "C:\\Program Files\\MATLAB\\R2024b",
		},
		{
			name:     "Unix MATLAB path",
			args:     []string{"--matlab-root=/usr/local/MATLAB/R2024b"},
			expected: "/usr/local/MATLAB/R2024b",
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.PreferredLocalMATLABRoot()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_PreferredMATLABStartingDirectory_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Windows custom project path",
			args:     []string{"--initial-working-folder=D:\\MATLAB_Projects"},
			expected: "D:\\MATLAB_Projects",
		},
		{
			name:     "Unix custom project path",
			args:     []string{"--initial-working-folder=/opt/matlab_work"},
			expected: "/opt/matlab_work",
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.PreferredMATLABStartingDirectory()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogDirectory_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: "",
		},
		{
			name:     "Supplied log directory",
			args:     []string{"--log-folder=/logs"},
			expected: "/logs",
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.BaseDir()

			// Assert
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogLevel_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name     string
		args     []string
		expected entities.LogLevel
	}{
		{
			name:     "default value",
			args:     []string{},
			expected: entities.LogLevelInfo,
		},
		{
			name:     "debug level",
			args:     []string{"--log-level=debug"},
			expected: entities.LogLevelDebug,
		},
		{
			name:     "info level",
			args:     []string{"--log-level=info"},
			expected: entities.LogLevelInfo,
		},
		{
			name:     "warn level",
			args:     []string{"--log-level=warn"},
			expected: entities.LogLevelWarn,
		},
		{
			name:     "error level",
			args:     []string{"--log-level=error"},
			expected: entities.LogLevelError,
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			// Act
			result := cfg.LogLevel()

			// Assert
			require.NoError(t, err)
			assert.Equal(t, testConfig.expected, result)
		})
	}
}

func TestConfig_LogLevel_Invalid(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Args().
		Return(append([]string{"testprocess"}, "--log-level=invalid")).
		Once()

	// Act
	cfg, err := config.New(mockOSLayer)

	// Assert
	require.Errorf(t, err, "invalid log level")
	assert.Empty(t, cfg)
}

func TestConfig_LogLevel_EmptyIsInvalid(t *testing.T) {
	// Arrange
	mockOSLayer := &configmocks.MockOSLayer{}
	defer mockOSLayer.AssertExpectations(t)

	mockOSLayer.EXPECT().
		Args().
		Return(append([]string{"testprocess"}, "--log-level=")).
		Once()

	// Act
	cfg, err := config.New(mockOSLayer)

	// Assert
	require.Errorf(t, err, "invalid log level")
	assert.Empty(t, cfg)
}

func TestConfig_Log_HappyPath(t *testing.T) {
	testConfigs := []struct {
		name                string
		args                []string
		expectedLogMessage  string
		expectedConfigField map[string]any
	}{
		{
			name:               "default configuration",
			args:               []string{},
			expectedLogMessage: "Configuration state",
			expectedConfigField: map[string]any{
				"disable-telemetry":         false,
				"initial-working-folder":    "",
				"log-level":                 entities.LogLevelInfo,
				"matlab-root":               "",
				"use-single-matlab-session": true,
			},
		},
		{
			name:               "custom configuration",
			args:               []string{"--disable-telemetry", "--use-single-matlab-session=false", "--log-level=debug", "--initial-working-folder=/home/user", "--matlab-root=/home/matlab"},
			expectedLogMessage: "Configuration state",
			expectedConfigField: map[string]any{
				"disable-telemetry":         true,
				"initial-working-folder":    "/home/user",
				"log-level":                 entities.LogLevelDebug,
				"matlab-root":               "/home/matlab",
				"use-single-matlab-session": false,
			},
		},
	}

	for _, testConfig := range testConfigs {
		t.Run(testConfig.name, func(t *testing.T) {
			// Arrange
			mockOSLayer := &configmocks.MockOSLayer{}
			defer mockOSLayer.AssertExpectations(t)

			mockOSLayer.EXPECT().
				Args().
				Return(append([]string{"testprocess"}, testConfig.args...)).
				Once()

			cfg, err := config.New(mockOSLayer)
			require.NoError(t, err)

			testLogger := testutils.NewInspectableLogger()

			// Act
			cfg.RecordToLogger(testLogger)

			// Assert
			infoLogs := testLogger.InfoLogs()
			require.Len(t, infoLogs, 1)

			fields, found := infoLogs[testConfig.expectedLogMessage]
			require.True(t, found, "Expected log message not found")

			for expectedField, expectedValue := range testConfig.expectedConfigField {
				actualValue, exists := fields[expectedField]
				require.True(t, exists, "%s field not found in log", expectedField)
				assert.Equal(t, expectedValue, actualValue, "%s field has incorrect value", expectedField)
			}
		})
	}
}
