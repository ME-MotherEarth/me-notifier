package shared

import (
	"net/http"

	"github.com/ME-MotherEarth/me-notifier/data"
	"github.com/gin-gonic/gin"
)

// HTTPServerCloser defines the basic actions of starting and closing that a web server should be able to do
type HTTPServerCloser interface {
	Start()
	Close() error
	IsInterfaceNil() bool
}

// GroupHandler defines the actions needed to be performed by an gin API group
type GroupHandler interface {
	RegisterRoutes(ws gin.IRoutes)
	GetAdditionalMiddlewares() []gin.HandlerFunc
	IsInterfaceNil() bool
}

// FacadeHandler defines the behavior of a notifier base facade handler
type FacadeHandler interface {
	HandlePushEvents(events data.SaveBlockData)
	HandleRevertEvents(revertBlock data.RevertBlock)
	HandleFinalizedEvents(finalizedBlock data.FinalizedBlock)
	GetConnectorUserAndPass() (string, string)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	IsInterfaceNil() bool
}

// WebServerHandler defines the behavior of a web server
type WebServerHandler interface {
	Run() error
	Close() error
	IsInterfaceNil() bool
}
