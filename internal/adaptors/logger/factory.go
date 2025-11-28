// Copyright 2025 The MathWorks, Inc.

package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultGlobalLogLevel slog.Level = slog.LevelDebug

const (
	logFileName         = "server"
	watchdogLogFileName = "watchdog"
	logFileExt          = ".log"
)

type Config interface {
	LogLevel() entities.LogLevel
	WatchdogMode() bool
}

type Directory interface {
	BaseDir() string
	ID() string
}

type FilenameFactory interface {
	FilenameWithSuffix(fileName string, ext string, suffix string) string
}

type OSLayer interface {
	Create(name string) (osfacade.File, error)
}

type Factory struct {
	globalLoggerOnce     *sync.Once
	globalLogger         *slogLogger
	globalLoggerLogLevel slog.Level
	globalLoggerFile     entities.Writer
}

func NewFactory(
	config Config,
	directory Directory,
	filenameFactory FilenameFactory,
	osLayer OSLayer,
) (*Factory, error) {
	logLevel, err := parseLogLevel(config.LogLevel())
	if err != nil {
		return nil, err
	}

	baseDir := directory.BaseDir()
	id := directory.ID()

	logFileBase := filepath.Join(baseDir, logFileName)
	if config.WatchdogMode() {
		logFileBase = filepath.Join(baseDir, watchdogLogFileName)
	}

	logFilePath := filenameFactory.FilenameWithSuffix(logFileBase, logFileExt, id)

	logFile, err := osLayer.Create(logFilePath)
	if err != nil {
		return nil, err
	}

	return &Factory{
		globalLoggerOnce:     new(sync.Once),
		globalLoggerLogLevel: logLevel,
		globalLoggerFile:     logFile,
	}, nil
}

func (f *Factory) NewMCPSessionLogger(session *mcp.ServerSession) entities.Logger {
	// In MCP Server development, special care should be given to logging:
	//
	// https://modelcontextprotocol.io/quickstart/server#logging-in-mcp-servers
	//
	// In essence, you can't log to standard `stdout`, and while you may log to `stderr`, you should log to the client:
	//
	// https://modelcontextprotocol.io/specification/2025-06-18/server/utilities/logging
	sessionHandler := mcp.NewLoggingHandler(session, &mcp.LoggingHandlerOptions{})

	handler := slog.NewJSONHandler(f.globalLoggerFile, &slog.HandlerOptions{
		Level: f.globalLoggerLogLevel,
	})

	return &slogLogger{
		logger: slog.New(NewMultiHandler(sessionHandler, handler)),
	}
}

func (f *Factory) GetGlobalLogger() entities.Logger {
	// There are cases where we want to log, but wo don't have an MCP session yet.
	// In those cases, we must log to stderr, to not affect the stdio transport:
	//
	// https://modelcontextprotocol.io/docs/develop/build-server#best-practices
	f.globalLoggerOnce.Do(func() {
		multiWriter := io.MultiWriter(os.Stderr, f.globalLoggerFile)

		handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level: f.globalLoggerLogLevel,
		})
		f.globalLogger = &slogLogger{
			logger: slog.New(handler),
		}
	})
	return f.globalLogger
}

func parseLogLevel(logLevel entities.LogLevel) (slog.Level, error) {
	switch logLevel {
	case entities.LogLevelDebug:
		return slog.LevelDebug, nil
	case entities.LogLevelInfo:
		return slog.LevelInfo, nil
	case entities.LogLevelWarn:
		return slog.LevelWarn, nil
	case entities.LogLevelError:
		return slog.LevelError, nil
	default:
		return defaultGlobalLogLevel, fmt.Errorf("unknown log level: %s", logLevel)
	}
}
