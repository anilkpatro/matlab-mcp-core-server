// Copyright 2025 The MathWorks, Inc.

package config

import (
	"fmt"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/inputs/flags"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/spf13/pflag"
)

func setupFlags(flagSet *pflag.FlagSet) error {
	flagSet.Bool(flags.VersionMode, flags.VersionModeDefaultValue,
		flags.VersionDescription,
	)

	flagSet.Bool(flags.DisableTelemetry, flags.DisableTelemetryDefaultValue,
		flags.DisableTelemetryDescription,
	)

	flagSet.Bool(flags.UseSingleMATLABSession, flags.UseSingleMATLABSessionDefaultValue,
		flags.UseSingleMATLABSessionDescription,
	)

	flagSet.String(flags.LogLevel, flags.LogLevelDefaultValue,
		flags.LogLevelDescription,
	)

	flagSet.String(flags.PreferredLocalMATLABRoot, flags.PreferredLocalMATLABRootDefaultValue,
		fmt.Sprintf("When %s is true, if this is set, defines which local MATLAB installation to use. If not set, the first MATLAB installation on the PATH will be used.", flags.UseSingleMATLABSession),
	)

	flagSet.String(flags.PreferredMATLABStartingDirectory, flags.PreferredMATLABStartingDirectoryDefaultValue,
		fmt.Sprintf("When %s is true, if this is set, defines which startup folder MATLAB will use. If not set, MATLAB will use the default MATLAB's startup folder.", flags.UseSingleMATLABSession),
	)

	flagSet.String(flags.BaseDir, flags.BaseDirDefaultValue,
		flags.BaseDirDescription,
	)

	// Hidden flags, for internal use only
	flagSet.Bool(flags.WatchdogMode, flags.WatchdogModeDefaultValue,
		flags.WatchdogModeDescription,
	)
	err := flagSet.MarkHidden(flags.WatchdogMode)
	if err != nil {
		return err
	}

	flagSet.String(flags.ServerInstanceID, flags.ServerInstanceIDDefaultValue,
		flags.ServerInstanceIDDescription,
	)
	err = flagSet.MarkHidden(flags.ServerInstanceID)
	if err != nil {
		return err
	}

	return nil
}

func createConfigWithFlagValues(osLayer OSLayer, flagSet *pflag.FlagSet, args []string) (*Config, error) {
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	versionMode, err := flagSet.GetBool(flags.VersionMode)
	if err != nil {
		return nil, err
	}

	disableTelemetry, err := flagSet.GetBool(flags.DisableTelemetry)
	if err != nil {
		return nil, err
	}

	useSingleMATLABSession, err := flagSet.GetBool(flags.UseSingleMATLABSession)
	if err != nil {
		return nil, err
	}

	logLevel, err := flagSet.GetString(flags.LogLevel)
	if err != nil {
		return nil, err
	}

	switch logLevel {
	case string(entities.LogLevelDebug), string(entities.LogLevelInfo), string(entities.LogLevelWarn), string(entities.LogLevelError):
		break
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevel)
	}

	preferredLocalMATLABRoot, err := flagSet.GetString(flags.PreferredLocalMATLABRoot)
	if err != nil {
		return nil, err
	}

	preferredMATLABStartingDirectory, err := flagSet.GetString(flags.PreferredMATLABStartingDirectory)
	if err != nil {
		return nil, err
	}

	baseDir, err := flagSet.GetString(flags.BaseDir)
	if err != nil {
		return nil, err
	}

	watchdogMode, err := flagSet.GetBool(flags.WatchdogMode)
	if err != nil {
		return nil, err
	}

	serverInstanceID, err := flagSet.GetString(flags.ServerInstanceID)
	if err != nil {
		return nil, err
	}

	return &Config{
		osLayer: osLayer,

		versionMode:                      versionMode,
		disableTelemetry:                 disableTelemetry,
		useSingleMATLABSession:           useSingleMATLABSession,
		logLevel:                         entities.LogLevel(logLevel),
		preferredLocalMATLABRoot:         preferredLocalMATLABRoot,
		preferredMATLABStartingDirectory: preferredMATLABStartingDirectory,
		baseDirectory:                    baseDir,
		watchdogMode:                     watchdogMode,
		serverInstanceID:                 serverInstanceID,
	}, nil
}
