package routes

import (
	"animoshi-api-go/src/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupWaifuRoutes(e *echo.Echo, client *mongo.Client) {
	// GET ROUTES
	e.GET("/waifu", func(c echo.Context) error {
		return lib.GetWaifu(c, client)
	})

	e.GET("/waifus", func(c echo.Context) error {
		return lib.GetWaifus(c, client)
	})
}
