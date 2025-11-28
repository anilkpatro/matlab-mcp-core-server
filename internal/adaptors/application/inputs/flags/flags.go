// Copyright 2025 The MathWorks, Inc.

package flags

const (
	VersionMode             = "version"
	VersionModeDefaultValue = false
	VersionDescription      = "Display the version of the MATLAB MCP Core Server."

	DisableTelemetry             = "disable-telemetry"
	DisableTelemetryDefaultValue = false
	DisableTelemetryDescription  = "Disable collection of usage data. By default, this software may collect information about you and your usage and send it to MathWorks. This data helps us improve our products and services."

	UseSingleMATLABSession             = "use-single-matlab-session"
	UseSingleMATLABSessionDefaultValue = true
	UseSingleMATLABSessionDescription  = "When true, a MATLAB session is started when a MATLAB MCP Core Server starts, and stopped when the server is shut down. When false, the server can manage multiple MATLAB sessions."

	PreferredLocalMATLABRoot             = "matlab-root"
	PreferredLocalMATLABRootDefaultValue = ""
	PreferredLocalMATLABRootDescription  = "The path to the MATLAB installation to use. If not specified, the server will use the first MATLAB installation it finds."

	PreferredMATLABStartingDirectory             = "initial-working-folder"
	PreferredMATLABStartingDirectoryDefaultValue = ""
	PreferredMATLABStartingDirectoryDescription  = "The directory to use as the initial working directory for MATLAB sessions. If not specified, the server will use the current working directory."

	BaseDir             = "log-folder"
	BaseDirDefaultValue = ""
	BaseDirDescription  = "The directory to use for log files. If not specified, the server will use the current working directory."

	LogLevel             = "log-level"
	LogLevelDefaultValue = "info"
	LogLevelDescription  = "The log level to use. Valid values are 'debug', 'info', 'warn', 'error', and 'fatal'."

	// Hidden

	WatchdogMode             = "watchdog"
	WatchdogModeDefaultValue = false
	WatchdogModeDescription  = "INTERNAL USE ONLY."

	ServerInstanceID             = "server-instance-id"
	ServerInstanceIDDefaultValue = ""
	ServerInstanceIDDescription  = "INTERNAL USE ONLY."
)
