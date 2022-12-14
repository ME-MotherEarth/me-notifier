package integrationTests

import (
	"github.com/ME-MotherEarth/me-notifier/config"
	"github.com/ME-MotherEarth/me-notifier/disabled"
	"github.com/ME-MotherEarth/me-notifier/dispatcher"
	"github.com/ME-MotherEarth/me-notifier/dispatcher/hub"
	"github.com/ME-MotherEarth/me-notifier/dispatcher/ws"
	"github.com/ME-MotherEarth/me-notifier/facade"
	"github.com/ME-MotherEarth/me-notifier/filters"
	"github.com/ME-MotherEarth/me-notifier/mocks"
	"github.com/ME-MotherEarth/me-notifier/process"
	"github.com/ME-MotherEarth/me-notifier/rabbitmq"
	"github.com/ME-MotherEarth/me-notifier/redis"
)

type testNotifier struct {
	Facade         FacadeHandler
	Publisher      PublisherHandler
	WSHandler      dispatcher.WSHandler
	RedisClient    *mocks.RedisClientMock
	RabbitMQClient *mocks.RabbitClientMock
}

// NewTestNotifierWithWS will create a notifier instance for websockets flow
func NewTestNotifierWithWS(cfg *config.GeneralConfig) (*testNotifier, error) {
	redisClient := mocks.NewRedisClientMock()
	redlockArgs := redis.ArgsRedlockWrapper{
		Client:       redisClient,
		TTLInMinutes: cfg.Redis.TTL,
	}
	locker, err := redis.NewRedlockWrapper(redlockArgs)
	if err != nil {
		return nil, err
	}

	args := hub.ArgsCommonHub{
		Filter:             filters.NewDefaultFilter(),
		SubscriptionMapper: dispatcher.NewSubscriptionMapper(),
	}
	publisher, err := hub.NewCommonHub(args)
	if err != nil {
		return nil, err
	}

	argsEventsHandler := process.ArgsEventsHandler{
		Config:    cfg.ConnectorApi,
		Locker:    locker,
		Publisher: publisher,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	upgrader, err := ws.NewWSUpgraderWrapper(1024, 1024)
	if err != nil {
		return nil, err
	}
	wsHandlerArgs := ws.ArgsWebSocketProcessor{
		Hub:      publisher,
		Upgrader: upgrader,
	}
	wsHandler, err := ws.NewWebSocketProcessor(wsHandlerArgs)
	if err != nil {
		return nil, err
	}

	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     cfg.ConnectorApi,
		WSHandler:     wsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: mocks.NewRabbitClientMock(),
	}, nil
}

// NewTestNotifierWithRabbitMq will create a notifier instance with rabbitmq
func NewTestNotifierWithRabbitMq(cfg *config.GeneralConfig) (*testNotifier, error) {
	redisClient := mocks.NewRedisClientMock()
	redlockArgs := redis.ArgsRedlockWrapper{
		Client:       redisClient,
		TTLInMinutes: cfg.Redis.TTL,
	}
	locker, err := redis.NewRedlockWrapper(redlockArgs)
	if err != nil {
		return nil, err
	}

	rabbitmqMock := mocks.NewRabbitClientMock()
	publisherArgs := rabbitmq.ArgsRabbitMqPublisher{
		Client: rabbitmqMock,
		Config: cfg.RabbitMQ,
	}
	publisher, err := rabbitmq.NewRabbitMqPublisher(publisherArgs)

	argsEventsHandler := process.ArgsEventsHandler{
		Config:    cfg.ConnectorApi,
		Locker:    locker,
		Publisher: publisher,
	}
	eventsHandler, err := process.NewEventsHandler(argsEventsHandler)
	if err != nil {
		return nil, err
	}

	wsHandler := &disabled.WSHandler{}
	facadeArgs := facade.ArgsNotifierFacade{
		EventsHandler: eventsHandler,
		APIConfig:     cfg.ConnectorApi,
		WSHandler:     wsHandler,
	}
	facade, err := facade.NewNotifierFacade(facadeArgs)
	if err != nil {
		return nil, err
	}

	return &testNotifier{
		Facade:         facade,
		Publisher:      publisher,
		WSHandler:      wsHandler,
		RedisClient:    redisClient,
		RabbitMQClient: rabbitmqMock,
	}, nil
}

func GetDefaultConfigs() *config.GeneralConfig {
	return &config.GeneralConfig{
		ConnectorApi: config.ConnectorApiConfig{
			Port:            "8081",
			Username:        "user",
			Password:        "pass",
			CheckDuplicates: false,
		},
		Redis: config.RedisConfig{
			Url:            "redis://localhost:6379",
			Channel:        "pub-sub",
			MasterName:     "mymaster",
			SentinelUrl:    "localhost:26379",
			ConnectionType: "sentinel",
			TTL:            30,
		},
		RabbitMQ: config.RabbitMQConfig{
			Url: "amqp://guest:guest@localhost:5672",
			EventsExchange: config.RabbitMQExchangeConfig{
				Name: "allevents",
				Type: "fanout",
			},
			RevertEventsExchange: config.RabbitMQExchangeConfig{
				Name: "revert",
				Type: "fanout",
			},
			FinalizedEventsExchange: config.RabbitMQExchangeConfig{
				Name: "finalized",
				Type: "fanout",
			},
			BlockTxsExchange: config.RabbitMQExchangeConfig{
				Name: "blocktxs",
				Type: "fanout",
			},
			BlockScrsExchange: config.RabbitMQExchangeConfig{
				Name: "blockscrs",
				Type: "fanout",
			},
		},
		Flags: &config.FlagsConfig{
			LogLevel:          "*:INFO",
			SaveLogFile:       false,
			GeneralConfigPath: "./config/config.toml",
			WorkingDir:        "",
			APIType:           "notifier",
		},
	}
}
