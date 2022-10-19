package factory

import (
	"github.com/ME-MotherEarth/me-notifier/common"
	"github.com/ME-MotherEarth/me-notifier/config"
	"github.com/ME-MotherEarth/me-notifier/process"
)

// ArgsEventsHandlerFactory defines the arguments needed for events handler creation
type ArgsEventsHandlerFactory struct {
	APIConfig    config.ConnectorApiConfig
	Locker       process.LockService
	MqPublisher  process.Publisher
	HubPublisher process.Publisher
	APIType      string
}

// CreateEventsHandler will create an events handler processor
func CreateEventsHandler(args ArgsEventsHandlerFactory) (process.EventsHandler, error) {
	publisher, err := getPublisher(args.APIType, args.MqPublisher, args.HubPublisher)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Config:    args.APIConfig,
		Locker:    args.Locker,
		Publisher: publisher,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	return eventsHandler, nil
}

func getPublisher(
	apiType string,
	mqPublisher process.Publisher,
	hubPublisher process.Publisher,
) (process.Publisher, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return mqPublisher, nil
	case common.WSAPIType:
		return hubPublisher, nil
	default:
		return nil, common.ErrInvalidAPIType
	}
}
