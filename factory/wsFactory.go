package factory

import (
	"github.com/ME-MotherEarth/me-notifier/common"
	"github.com/ME-MotherEarth/me-notifier/disabled"
	"github.com/ME-MotherEarth/me-notifier/dispatcher"
	"github.com/ME-MotherEarth/me-notifier/dispatcher/ws"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

// CreateWSHandler creates websocket handler component based on api type
func CreateWSHandler(apiType string, hub dispatcher.Hub) (dispatcher.WSHandler, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.WSHandler{}, nil
	case common.WSAPIType:
		return createWSHandler(hub)
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createWSHandler(hub dispatcher.Hub) (dispatcher.WSHandler, error) {
	upgrader, err := ws.NewWSUpgraderWrapper(readBufferSize, writeBufferSize)
	if err != nil {
		return nil, err
	}

	args := ws.ArgsWebSocketProcessor{
		Hub:      hub,
		Upgrader: upgrader,
	}
	return ws.NewWebSocketProcessor(args)
}
