package mocks

import (
	"net/http"

	"github.com/ME-MotherEarth/me-notifier/data"
)

// FacadeStub implements FacadeHandler interface
type FacadeStub struct {
	HandlePushEventsCalled        func(events data.SaveBlockData)
	HandleRevertEventsCalled      func(events data.RevertBlock)
	HandleFinalizedEventsCalled   func(events data.FinalizedBlock)
	ServeCalled                   func(w http.ResponseWriter, r *http.Request)
	GetConnectorUserAndPassCalled func() (string, string)
}

// HandlePushEvents -
func (fs *FacadeStub) HandlePushEvents(events data.SaveBlockData) {
	if fs.HandlePushEventsCalled != nil {
		fs.HandlePushEventsCalled(events)
	}
}

// HandleRevertEvents -
func (fs *FacadeStub) HandleRevertEvents(events data.RevertBlock) {
	if fs.HandleRevertEventsCalled != nil {
		fs.HandleRevertEventsCalled(events)
	}
}

// HandleFinalizedEvents -
func (fs *FacadeStub) HandleFinalizedEvents(events data.FinalizedBlock) {
	if fs.HandleFinalizedEventsCalled != nil {
		fs.HandleFinalizedEventsCalled(events)
	}
}

// ServeHTTP -
func (fs *FacadeStub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if fs.ServeCalled != nil {
		fs.ServeCalled(w, r)
	}
}

// GetConnectorUserAndPass -
func (fs *FacadeStub) GetConnectorUserAndPass() (string, string) {
	if fs.GetConnectorUserAndPassCalled != nil {
		return fs.GetConnectorUserAndPassCalled()
	}

	return "", ""
}

// IsInterfaceNil -
func (fs *FacadeStub) IsInterfaceNil() bool {
	return fs == nil
}
