package factory

import (
	"github.com/ME-MotherEarth/me-notifier/common"
	"github.com/ME-MotherEarth/me-notifier/disabled"
	"github.com/ME-MotherEarth/me-notifier/dispatcher"
	"github.com/ME-MotherEarth/me-notifier/dispatcher/hub"
	"github.com/ME-MotherEarth/me-notifier/filters"
)

// CreateHub creates a common hub component
func CreateHub(apiType string) (dispatcher.Hub, error) {
	switch apiType {
	case common.MessageQueueAPIType:
		return &disabled.Hub{}, nil
	case common.WSAPIType:
		return createHub()
	default:
		return nil, common.ErrInvalidAPIType
	}
}

func createHub() (dispatcher.Hub, error) {
	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	return hub.NewCommonHub(args)
}
