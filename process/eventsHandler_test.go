package process_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ME-MotherEarth/me-core/core/check"
	"github.com/ME-MotherEarth/me-core/data/smartContractResult"
	"github.com/ME-MotherEarth/me-core/data/transaction"
	"github.com/ME-MotherEarth/me-notifier/config"
	"github.com/ME-MotherEarth/me-notifier/data"
	"github.com/ME-MotherEarth/me-notifier/mocks"
	"github.com/ME-MotherEarth/me-notifier/process"
	"github.com/stretchr/testify/require"
)

func createMockEventsHandlerArgs() process.ArgsEventsHandler {
	return process.ArgsEventsHandler{
		Config: config.ConnectorApiConfig{
			CheckDuplicates: false,
		},
		Locker:    &mocks.LockerStub{},
		Publisher: &mocks.PublisherStub{},
	}
}

func TestNewEventsHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil locker service", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Locker = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, process.ErrNilLockService, err)
		require.True(t, check.IfNil(eventsHandler))
	})

	t.Run("nil publisher", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Publisher = nil

		eventsHandler, err := process.NewEventsHandler(args)
		require.Equal(t, process.ErrNilPublisherService, err)
		require.Nil(t, eventsHandler)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)
		require.NotNil(t, eventsHandler)
	})
}

func TestHandlePushEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastCalled: func(events data.BlockEvents) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockEvents{
			Hash:   "hash1",
			Events: []data.Event{},
		}

		eventsHandler.HandlePushEvents(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastCalled: func(events data.BlockEvents) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockEvents{
			Hash:   "hash1",
			Events: []data.Event{},
		}

		eventsHandler.HandlePushEvents(events)
		require.False(t, wasCalled)
	})
}

func TestHandleRevertEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastRevertCalled: func(events data.RevertBlock) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.RevertBlock{
			Hash:  "hash1",
			Nonce: 1,
		}

		eventsHandler.HandleRevertEvents(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastRevertCalled: func(events data.RevertBlock) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.RevertBlock{
			Hash:  "hash1",
			Nonce: 1,
		}

		eventsHandler.HandleRevertEvents(events)
		require.False(t, wasCalled)
	})
}

func TestHandleFinalizedEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast finalized event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastFinalizedCalled: func(events data.FinalizedBlock) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.FinalizedBlock{
			Hash: "hash1",
		}

		eventsHandler.HandleFinalizedEvents(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastFinalizedCalled: func(events data.FinalizedBlock) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.FinalizedBlock{
			Hash: "hash1",
		}

		eventsHandler.HandleFinalizedEvents(events)
		require.False(t, wasCalled)
	})
}

func TestHandleTxsEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast txs event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastTxsCalled: func(event data.BlockTxs) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockTxs{
			Hash: "hash1",
			Txs: map[string]transaction.Transaction{
				"hash1": {
					Nonce: 1,
				},
			},
		}

		eventsHandler.HandleBlockTxs(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastTxsCalled: func(event data.BlockTxs) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockTxs{
			Hash: "hash1",
			Txs: map[string]transaction.Transaction{
				"hash1": {
					Nonce: 1,
				},
			},
		}

		eventsHandler.HandleBlockTxs(events)
		require.False(t, wasCalled)
	})
}

func TestHandleScrsEvents(t *testing.T) {
	t.Parallel()

	t.Run("broadcast scrs event was called", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Publisher = &mocks.PublisherStub{
			BroadcastScrsCalled: func(event data.BlockScrs) {
				wasCalled = true
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockScrs{
			Hash: "hash1",
			Scrs: map[string]smartContractResult.SmartContractResult{
				"hash2": {
					Nonce: 2,
				},
			},
		}

		eventsHandler.HandleBlockScrs(events)
		require.True(t, wasCalled)
	})

	t.Run("check duplicates enabled, should not process event", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Publisher = &mocks.PublisherStub{
			BroadcastScrsCalled: func(event data.BlockScrs) {
				wasCalled = true
			},
		}
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		events := data.BlockScrs{
			Hash: "hash1",
			Scrs: map[string]smartContractResult.SmartContractResult{
				"hash2": {
					Nonce: 2,
				},
			},
		}

		eventsHandler.HandleBlockScrs(events)
		require.False(t, wasCalled)
	})
}

func TestTryCheckProcessedWithRetry(t *testing.T) {
	t.Parallel()

	hash := "hash1"

	t.Run("event is NOT already processed", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return false, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(hash)
		require.False(t, ok)
	})

	t.Run("event is already processed", func(t *testing.T) {
		t.Parallel()

		args := createMockEventsHandlerArgs()
		args.Config.CheckDuplicates = true
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				return true, nil
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(hash)
		require.True(t, ok)
	})

	t.Run("locker service is failing on first try, has no connection, works on second try", func(t *testing.T) {
		t.Parallel()

		numCallsHasConnection := 0
		numCallsIsProcessed := 0

		args := createMockEventsHandlerArgs()
		args.Locker = &mocks.LockerStub{
			IsEventProcessedCalled: func(ctx context.Context, blockHash string) (bool, error) {
				if numCallsIsProcessed > 0 {
					return true, nil
				}
				numCallsIsProcessed++
				return false, errors.New("fail to process")
			},
			HasConnectionCalled: func(ctx context.Context) bool {
				if numCallsHasConnection > 0 {
					return true
				}
				numCallsHasConnection++
				return false
			},
		}

		eventsHandler, err := process.NewEventsHandler(args)
		require.Nil(t, err)

		ok := eventsHandler.TryCheckProcessedWithRetry(hash)
		require.True(t, ok)
	})
}
