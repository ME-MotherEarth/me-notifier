package integrationTests

import (
	"net/http"

	"github.com/ME-MotherEarth/me-notifier/data"
)

// FacadeHandler defines facade behaviour
type FacadeHandler interface {
	HandlePushEvents(events data.SaveBlockData)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}

// PublisherHandler defines publisher behaviour
type PublisherHandler interface {
	Run()
	Broadcast(events data.BlockEvents)
	BroadcastRevert(event data.RevertBlock)
	BroadcastFinalized(event data.FinalizedBlock)
	Close() error
	IsInterfaceNil() bool
}
