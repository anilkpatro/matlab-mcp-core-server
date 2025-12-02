// Copyright 2025 The MathWorks, Inc.

package globalmatlab_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab"
	entitiesmocks "github.com/matlab/matlab-mcp-core-server/mocks/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGlobalMATLAB_Client_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedPreferredMATLABStartingDir := filepath.Join("some", "starting", "dir")

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedPreferredMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedPreferredMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_StartingDirectorySet(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedPreferredMATLABStartingDir := filepath.Join("some", "starting", "dir")

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedPreferredMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedPreferredMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_SelectFirstMATLABVersionOnPathError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return("", expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_MATLABStartingDirSelectionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	expectedError := assert.AnError

	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: false,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return("", expectedError).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	require.NotNil(t, globalMATLABSession)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedSessionClient, client)
}

func TestGlobalMATLAB_Client_StartMATLABSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	const sessionIDThatShouldBeUnused = entities.SessionID(0)
	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedPreferredMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedPreferredMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedPreferredMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionIDThatShouldBeUnused, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_GetMATLABSessionClientError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	sessionID := entities.SessionID(123)
	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedPreferredMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedPreferredMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedPreferredMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), sessionID).
		Return(nil, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client, err := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	assert.Nil(t, client)
}

func TestGlobalMATLAB_Client_ReturnsInitializeCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	expectedError := assert.AnError

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return("", expectedError).
		Once()

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	assert.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_ReturnsMATLABStartupCachedErrorOnSubsequentClientCalls(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	ctx := t.Context()
	const sessionIDThatShouldBeUnused = entities.SessionID(0)
	expectedPreferredMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedPreferredMATLABStartingDir := filepath.Join("some", "starting", "dir")
	expectedError := assert.AnError

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedPreferredMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedPreferredMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedPreferredMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedPreferredMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.Anything, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Return(sessionIDThatShouldBeUnused, expectedError).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	client1, err1 := globalMATLABSession.Client(ctx, mockLogger)
	client2, err2 := globalMATLABSession.Client(ctx, mockLogger)

	// Assert
	assert.Nil(t, client1)
	require.ErrorIs(t, err1, expectedError)

	assert.Nil(t, client2)
	require.ErrorIs(t, err2, expectedError)
}

func TestGlobalMATLAB_Client_ConcurrentCallsWaitForCompletion(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	expectedSessionClient := &entitiesmocks.MockMATLABSessionClient{}

	ctx := t.Context()
	expectedSelectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	mockSessionID := entities.SessionID(123)

	expectedLocalSessionDetails := entities.LocalSessionDetails{
		MATLABRoot:             expectedSelectedMATLABRoot,
		IsStartingDirectorySet: true,
		StartingDirectory:      expectedMATLABStartingDir,
		ShowMATLABDesktop:      true,
	}

	blockStartMATLAB := make(chan struct{})
	startMATLABCalled := make(chan struct{})

	type clientResult struct {
		client entities.MATLABSessionClient
		err    error
	}

	firstCallCompleted := make(chan clientResult)
	secondCallCompleted := make(chan clientResult)

	mockMATLABRootSelector.EXPECT().
		SelectFirstMATLABVersionOnPath(ctx, mockLogger.AsMockArg()).
		Return(expectedSelectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMatlabStartingDir().
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(ctx, mockLogger.AsMockArg(), expectedLocalSessionDetails).
		Run(func(ctx context.Context, logger entities.Logger, details entities.SessionDetails) {
			close(startMATLABCalled)
			<-blockStartMATLAB
		}).
		Return(mockSessionID, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), mockSessionID).
		Return(expectedSessionClient, nil).
		Once()

	mockMATLABManager.EXPECT().
		GetMATLABSessionClient(ctx, mockLogger.AsMockArg(), mockSessionID).
		Return(expectedSessionClient, nil).
		Once()

	globalMATLABSession := globalmatlab.New(
		mockMATLABManager,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		firstCallCompleted <- clientResult{client: client, err: err}
	}()

	<-startMATLABCalled

	go func() {
		client, err := globalMATLABSession.Client(ctx, mockLogger)
		secondCallCompleted <- clientResult{client: client, err: err}
	}()

	select {
	case <-secondCallCompleted:
		t.Fatal("Second Client call completed before first call was unblocked")
	case <-time.After(100 * time.Millisecond):
		// Second call is still waiting
	}

	close(blockStartMATLAB)
	result1 := <-firstCallCompleted
	result2 := <-secondCallCompleted

	// Assert
	require.NoError(t, result1.err)
	assert.Equal(t, expectedSessionClient, result1.client)

	require.NoError(t, result2.err)
	assert.Equal(t, expectedSessionClient, result2.client)
}
