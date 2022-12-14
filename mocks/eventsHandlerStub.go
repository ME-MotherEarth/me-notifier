package mocks

import "github.com/ME-MotherEarth/me-notifier/data"

// EventsHandlerStub implements EventsHandler interface
type EventsHandlerStub struct {
	HandlePushEventsCalled      func(events data.BlockEvents)
	HandleRevertEventsCalled    func(revertBlock data.RevertBlock)
	HandleFinalizedEventsCalled func(finalizedBlock data.FinalizedBlock)
	HandleBlockTxsCalled        func(blockTxs data.BlockTxs)
	HandleBlockScrsCalled       func(blockScrs data.BlockScrs)
}

// HandlePushEvents -
func (e *EventsHandlerStub) HandlePushEvents(events data.BlockEvents) {
	if e.HandlePushEventsCalled != nil {
		e.HandlePushEventsCalled(events)
	}
}

// HandleRevertEvents -
func (e *EventsHandlerStub) HandleRevertEvents(revertBlock data.RevertBlock) {
	if e.HandleRevertEventsCalled != nil {
		e.HandleRevertEventsCalled(revertBlock)
	}
}

// HandleFinalizedEvents -
func (e *EventsHandlerStub) HandleFinalizedEvents(finalizedBlock data.FinalizedBlock) {
	if e.HandleFinalizedEventsCalled != nil {
		e.HandleFinalizedEventsCalled(finalizedBlock)
	}
}

// HandleBlockTxs -
func (e *EventsHandlerStub) HandleBlockTxs(blockTxs data.BlockTxs) {
	if e.HandleBlockTxsCalled != nil {
		e.HandleBlockTxsCalled(blockTxs)
	}
}

// HandleBlockScrs -
func (e *EventsHandlerStub) HandleBlockScrs(blockScrs data.BlockScrs) {
	if e.HandleBlockScrsCalled != nil {
		e.HandleBlockScrsCalled(blockScrs)
	}
}

// IsInterfaceNil -
func (e *EventsHandlerStub) IsInterfaceNil() bool {
	return e == nil
}
