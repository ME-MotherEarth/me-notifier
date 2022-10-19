package gin_test

import (
	"testing"

	"github.com/ME-MotherEarth/me-core/core/check"
	apiErrors "github.com/ME-MotherEarth/me-notifier/api/errors"
	"github.com/ME-MotherEarth/me-notifier/api/gin"
	"github.com/ME-MotherEarth/me-notifier/common"
	"github.com/ME-MotherEarth/me-notifier/config"
	"github.com/ME-MotherEarth/me-notifier/mocks"
	"github.com/stretchr/testify/require"
)

func createMockArgsWebServerHandler() gin.ArgsWebServerHandler {
	return gin.ArgsWebServerHandler{
		Facade: &mocks.FacadeStub{},
		Config: config.ConnectorApiConfig{
			Port: "8080",
		},
		Type: "notifier",
	}
}

func TestNewWebServerHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil facade", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Facade = nil

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, apiErrors.ErrNilFacadeHandler, err)
	})

	t.Run("invalid api type", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()
		args.Type = ""

		ws, err := gin.NewWebServerHandler(args)
		require.True(t, check.IfNil(ws))
		require.Equal(t, common.ErrInvalidAPIType, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createMockArgsWebServerHandler()

		ws, err := gin.NewWebServerHandler(args)
		require.Nil(t, err)
		require.NotNil(t, ws)

		err = ws.Run()
		require.Nil(t, err)

		err = ws.Close()
		require.Nil(t, err)
	})
}
