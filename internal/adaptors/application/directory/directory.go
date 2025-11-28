// Copyright 2025 The MathWorks, Inc.

package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
)

const (
	defaultLogDirPattern = "matlab-mcp-core-server-"
	markerFileName       = ".matlab-mcp-core-server"
)

type Config interface {
	BaseDir() string
	ServerInstanceID() string
}

type FilenameFactory interface {
	CreateFileWithUniqueSuffix(baseName string, ext string) (string, string, error)
}

type OSLayer interface {
	MkdirTemp(dir string, pattern string) (string, error)
	MkdirAll(name string, perm os.FileMode) error
	Create(name string) (osfacade.File, error)
}

type Directory struct {
	baseDir string
	id      string

	osFacade OSLayer
}

func New(
	config Config,
	filenameFactory FilenameFactory,
	osFacade OSLayer,
) (*Directory, error) {
	baseDir := config.BaseDir()

	if baseDir == "" {
		var err error
		if baseDir, err = osFacade.MkdirTemp("", defaultLogDirPattern); err != nil {
			return nil, err
		}
	} else {
		if err := osFacade.MkdirAll(baseDir, 0o700); err != nil {
			return nil, err
		}
	}

	serverInstanceID := config.ServerInstanceID()

	if serverInstanceID == "" {
		_, id, err := filenameFactory.CreateFileWithUniqueSuffix(filepath.Join(baseDir, markerFileName), "")
		if err != nil {
			return nil, err
		}

		serverInstanceID = id
	}

	return &Directory{
		baseDir: baseDir,
		id:      serverInstanceID,

		osFacade: osFacade,
	}, nil
}

func (d *Directory) BaseDir() string {
	return d.baseDir
}

func (d *Directory) ID() string {
	return d.id
}

func (d *Directory) CreateSubDir(pattern string) (string, error) {
	if !strings.HasSuffix(pattern, "-") {
		pattern = fmt.Sprintf("%s-", pattern)
	}

	pattern = fmt.Sprintf("%s%s-", pattern, d.id)

	return d.osFacade.MkdirTemp(d.baseDir, pattern)
}

func (d *Directory) RecordToLogger(logger entities.Logger) {
	logger.
		With("log-dir", d.baseDir).
		With("id", d.id).
		Info("Application directory state")
}
