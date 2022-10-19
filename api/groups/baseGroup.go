package groups

import (
	logger "github.com/ME-MotherEarth/me-logger"
	"github.com/ME-MotherEarth/me-notifier/api/shared"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints []*shared.EndpointHandlerData
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(
	ws gin.IRoutes,
) {
	for _, handlerData := range bg.endpoints {
		ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
	}
}
