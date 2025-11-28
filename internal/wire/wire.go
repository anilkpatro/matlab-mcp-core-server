// Copyright 2025 The MathWorks, Inc.

//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/directory"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/lifecyclesignaler"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/modeselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/orchestrator"
	files "github.com/matlab/matlab-mcp-core-server/internal/adaptors/filesystem/files"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/matlabrootselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/matlabstartingdirselector"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/logger"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/directorymanager/matlabfiles"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processdetails"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/localmatlabsession/processlauncher"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabroot"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabservices/services/matlablocator/matlabversion"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionclient"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/matlabmanager/matlabsessionstore"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/server/configurator"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/basetool"
	evalmatlabcodemultisessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/evalmatlabcode"
	listavailablematlabstool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/listavailablematlabs"
	startmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/startmatlabsession"
	stopmatlabsessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/multisession/stopmatlabsession"
	checkmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/checkmatlabcode"
	detectmatlabtoolboxessinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/detectmatlabtoolboxes"
	evalmatlabcodesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/evalmatlabcode"
	runmatlabfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabfile"
	runmatlabtestfilesinglesessiontool "github.com/matlab/matlab-mcp-core-server/internal/adaptors/mcp/tools/singlesession/runmatlabtestfile"
	watchdogclient "github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/watchdog/process"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/filefacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/iofacade"
	"github.com/matlab/matlab-mcp-core-server/internal/facades/osfacade"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/checkmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/detectmatlabtoolboxes"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/evalmatlabcode"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/listavailablematlabs"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/runmatlabtestfile"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/startmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/stopmatlabsession"
	"github.com/matlab/matlab-mcp-core-server/internal/usecases/utils/pathvalidator"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/httpclientfactory"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/ossignaler"
	"github.com/matlab/matlab-mcp-core-server/internal/utils/oswrapper"
	watchdogprocess "github.com/matlab/matlab-mcp-core-server/internal/watchdog"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/processhandler"
	"github.com/matlab/matlab-mcp-core-server/internal/watchdog/transport"
)

type orchestratorFactory struct{}

func newOrchestratorFactory() *orchestratorFactory {
	return &orchestratorFactory{}
}

func (f *orchestratorFactory) Create() (entities.Mode, error) {
	return initializeOrchestrator()
}

type watchdogProcessFactory struct{}

func newWatchdogProcessFactory() *watchdogProcessFactory {
	return &watchdogProcessFactory{}
}

func (f *watchdogProcessFactory) Create() (entities.Mode, error) {
	return initializeWatchdog()
}

func InitializeModeSelector() (*modeselector.ModeSelector, error) {
	wire.Build(
		// Application
		modeselector.New,
		wire.Bind(new(modeselector.Config), new(*config.Config)),
		wire.Bind(new(modeselector.WatchdogProcessFactory), new(*watchdogProcessFactory)),
		wire.Bind(new(modeselector.OrchestratorFactory), new(*orchestratorFactory)),
		wire.Bind(new(modeselector.OSLayer), new(*osfacade.OsFacade)),

		// Factories
		newWatchdogProcessFactory,
		newOrchestratorFactory,

		// Low-level Interfaces
		config.New,
		wire.Bind(new(config.OSLayer), new(*osfacade.OsFacade)),
		osfacade.New,
	)

	return nil, nil
}

func initializeOrchestrator() (*orchestrator.Orchestrator, error) {
	wire.Build(
		// Orchestrator
		orchestrator.New,
		wire.Bind(new(orchestrator.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(orchestrator.Config), new(*config.Config)),
		wire.Bind(new(orchestrator.Server), new(*server.Server)),
		wire.Bind(new(orchestrator.WatchdogClient), new(*watchdogclient.Watchdog)),
		wire.Bind(new(orchestrator.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(orchestrator.OSSignaler), new(*ossignaler.OSSignaler)),
		wire.Bind(new(orchestrator.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(orchestrator.Directory), new(*directory.Directory)),

		// Watchdog Client
		watchdogclient.New,
		wire.Bind(new(watchdogclient.WatchdogProcess), new(*process.Process)),
		wire.Bind(new(watchdogclient.TransportFactory), new(*transport.Factory)),
		wire.Bind(new(watchdogclient.LoggerFactory), new(*logger.Factory)),

		// Watchdog Process Handler for Watchdog Client
		process.New,
		wire.Bind(new(process.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(process.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(process.Directory), new(*directory.Directory)),

		// Watchdog Transport Factory
		transport.NewFactory,

		// MCP Server
		server.NewMCPSDKServer,
		wire.Bind(new(server.ServerConfig), new(*config.Config)),
		server.New,
		wire.Bind(new(server.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(server.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),
		wire.Bind(new(server.MCPServerConfigurator), new(*configurator.Configurator)),

		// MCP Server Configurator
		configurator.New,
		wire.Bind(new(configurator.Config), new(*config.Config)),

		// Tools
		wire.Bind(new(basetool.LoggerFactory), new(*logger.Factory)),

		listavailablematlabstool.New,
		wire.Bind(new(listavailablematlabstool.Usecase), new(*listavailablematlabs.Usecase)),

		startmatlabsessiontool.New,
		wire.Bind(new(startmatlabsessiontool.Usecase), new(*startmatlabsession.Usecase)),

		stopmatlabsessiontool.New,
		wire.Bind(new(stopmatlabsessiontool.Usecase), new(*stopmatlabsession.Usecase)),

		evalmatlabcodemultisessiontool.New,
		wire.Bind(new(evalmatlabcodemultisessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		evalmatlabcodesinglesessiontool.New,
		wire.Bind(new(evalmatlabcodesinglesessiontool.Usecase), new(*evalmatlabcode.Usecase)),

		checkmatlabcodesinglesessiontool.New,
		wire.Bind(new(checkmatlabcodesinglesessiontool.Usecase), new(*checkmatlabcode.Usecase)),

		detectmatlabtoolboxessinglesessiontool.New,
		wire.Bind(new(detectmatlabtoolboxessinglesessiontool.Usecase), new(*detectmatlabtoolboxes.Usecase)),

		runmatlabfilesinglesessiontool.New,
		wire.Bind(new(runmatlabfilesinglesessiontool.Usecase), new(*runmatlabfile.Usecase)),

		runmatlabtestfilesinglesessiontool.New,
		wire.Bind(new(runmatlabtestfilesinglesessiontool.Usecase), new(*runmatlabtestfile.Usecase)),

		// Use Cases
		listavailablematlabs.New,
		startmatlabsession.New,
		stopmatlabsession.New,
		evalmatlabcode.New,
		wire.Bind(new(evalmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),
		checkmatlabcode.New,
		wire.Bind(new(checkmatlabcode.PathValidator), new(*pathvalidator.PathValidator)),
		detectmatlabtoolboxes.New,
		runmatlabfile.New,
		wire.Bind(new(runmatlabfile.PathValidator), new(*pathvalidator.PathValidator)),
		runmatlabtestfile.New,
		wire.Bind(new(runmatlabtestfile.PathValidator), new(*pathvalidator.PathValidator)),

		// Use Cases Utilities
		pathvalidator.New,
		wire.Bind(new(pathvalidator.OSLayer), new(*osfacade.OsFacade)),

		// Entities
		wire.Bind(new(entities.GlobalMATLAB), new(*globalmatlab.GlobalMATLAB)),
		wire.Bind(new(entities.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Manager
		matlabmanager.New,
		wire.Bind(new(matlabmanager.MATLABServices), new(*matlabservices.MATLABServices)),
		wire.Bind(new(matlabmanager.MATLABSessionStore), new(*matlabsessionstore.Store)),
		wire.Bind(new(matlabmanager.MATLABSessionClientFactory), new(*matlabsessionclient.Factory)),

		// MATLAB Session Store
		matlabsessionstore.New,
		wire.Bind(new(matlabsessionstore.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(matlabsessionstore.LifecycleSignaler), new(*lifecyclesignaler.LifecycleSignaler)),

		// MATLAB Session Client Factory
		matlabsessionclient.NewFactory,
		wire.Bind(new(matlabsessionclient.HttpClientFactory), new(*httpclientfactory.HTTPClientFactory)),

		// Global MATLAB Session
		globalmatlab.New,
		wire.Bind(new(globalmatlab.MATLABManager), new(*matlabmanager.MATLABManager)),
		wire.Bind(new(globalmatlab.MATLABRootSelector), new(*matlabrootselector.MATLABRootSelector)),
		wire.Bind(new(globalmatlab.MATLABStartingDirSelector), new(*matlabstartingdirselector.MATLABStartingDirSelector)),

		// MATLAB Root Selector
		matlabrootselector.New,
		wire.Bind(new(matlabrootselector.Config), new(*config.Config)),
		wire.Bind(new(matlabrootselector.MATLABManager), new(*matlabmanager.MATLABManager)),

		// MATLAB Starting Dir Selector
		matlabstartingdirselector.New,
		wire.Bind(new(matlabstartingdirselector.Config), new(*config.Config)),
		wire.Bind(new(matlabstartingdirselector.OSLayer), new(*osfacade.OsFacade)),

		// MATLAB Services
		matlabservices.New,
		wire.Bind(new(matlabservices.MATLABLocator), new(*matlablocator.MATLABLocator)),
		wire.Bind(new(matlabservices.LocalMATLABSessionLauncher), new(*localmatlabsession.Starter)),

		// MATLAB Locator
		matlablocator.New,
		wire.Bind(new(matlablocator.MATLABRootGetter), new(*matlabroot.Getter)),
		wire.Bind(new(matlablocator.MATLABVersionGetter), new(*matlabversion.Getter)),

		// Local MATLAB Session
		localmatlabsession.NewStarter,
		wire.Bind(new(localmatlabsession.SessionDirectoryFactory), new(*directorymanager.DirectoryFactory)),
		wire.Bind(new(localmatlabsession.ProcessDetails), new(*processdetails.ProcessDetails)),
		wire.Bind(new(localmatlabsession.MATLABProcessLauncher), new(*processlauncher.MATLABProcessLauncher)),
		wire.Bind(new(localmatlabsession.Watchdog), new(*watchdogclient.Watchdog)),

		// Local MATLAB Session Directory Manager
		directorymanager.NewFactory,
		wire.Bind(new(directorymanager.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(directorymanager.ApplicationDirectory), new(*directory.Directory)),
		wire.Bind(new(directorymanager.MATLABFiles), new(matlabfiles.MATLABFiles)),

		// Local MATLAB Session Process Details
		processdetails.New,
		wire.Bind(new(processdetails.OSLayer), new(*osfacade.OsFacade)),

		// Local MATLAB Process Launcher
		processlauncher.New,

		wire.NewSet(
			// MATLAB Root Getter
			matlabroot.New,
			wire.Bind(new(matlabroot.OSLayer), new(*osfacade.OsFacade)),
			wire.Bind(new(matlabroot.FileLayer), new(*filefacade.FileFacade)),

			// MATLAB Version Getter
			matlabversion.New,
			wire.Bind(new(matlabversion.OSLayer), new(*osfacade.OsFacade)),
			wire.Bind(new(matlabversion.IOLayer), new(*iofacade.IoFacade)),

			// MATLAB Files Provider
			matlabfiles.New,

			wire.NewSet(
				// Low-level Interfaces
				logger.NewFactory,
				wire.Bind(new(logger.Config), new(*config.Config)),
				wire.Bind(new(logger.Directory), new(*directory.Directory)),
				wire.Bind(new(logger.FilenameFactory), new(*files.Factory)),
				wire.Bind(new(logger.OSLayer), new(*osfacade.OsFacade)),
				oswrapper.New,
				wire.Bind(new(oswrapper.OSLayer), new(*osfacade.OsFacade)),
				directory.New,
				wire.Bind(new(directory.Config), new(*config.Config)),
				wire.Bind(new(directory.FilenameFactory), new(*files.Factory)),
				wire.Bind(new(directory.OSLayer), new(*osfacade.OsFacade)),
				lifecyclesignaler.New,
				config.New,
				wire.Bind(new(config.OSLayer), new(*osfacade.OsFacade)),
				files.NewFactory,
				wire.Bind(new(files.OSLayer), new(*osfacade.OsFacade)),
				osfacade.New,
				iofacade.New,
				filefacade.New,
				ossignaler.New,
				httpclientfactory.New,
			),
		),
	)

	return nil, nil
}

func initializeWatchdog() (*watchdogprocess.Watchdog, error) {
	wire.Build( // Watchdog Process
		watchdogprocess.New,
		wire.Bind(new(watchdogprocess.LoggerFactory), new(*logger.Factory)),
		wire.Bind(new(watchdogprocess.OSLayer), new(*osfacade.OsFacade)),
		wire.Bind(new(watchdogprocess.ProcessHandler), new(*processhandler.ProcessHandler)),
		wire.Bind(new(watchdogprocess.OSSignaler), new(*ossignaler.OSSignaler)),
		wire.Bind(new(watchdogprocess.TransportFactory), new(*transport.Factory)),

		// Process Handler for Watchdog Process
		processhandler.New,
		wire.Bind(new(processhandler.OSWrapper), new(*oswrapper.OSWrapper)),

		// Watchdog Transport Factory
		transport.NewFactory,

		// Low-level Interfaces
		logger.NewFactory,
		wire.Bind(new(logger.Config), new(*config.Config)),
		wire.Bind(new(logger.Directory), new(*directory.Directory)),
		wire.Bind(new(logger.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(logger.OSLayer), new(*osfacade.OsFacade)),
		directory.New,
		wire.Bind(new(directory.Config), new(*config.Config)),
		wire.Bind(new(directory.FilenameFactory), new(*files.Factory)),
		wire.Bind(new(directory.OSLayer), new(*osfacade.OsFacade)),
		config.New,
		wire.Bind(new(config.OSLayer), new(*osfacade.OsFacade)),
		files.NewFactory,
		wire.Bind(new(files.OSLayer), new(*osfacade.OsFacade)),
		oswrapper.New,
		wire.Bind(new(oswrapper.OSLayer), new(*osfacade.OsFacade)),
		ossignaler.New,
		osfacade.New,
	)

	return nil, nil
}
