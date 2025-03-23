package main

import (
	"animoshi-api-go/src/infra"
	"animoshi-api-go/src/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

var client *mongo.Client

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:4173", "http://localhost:3000", "https://animoshi-svelte-frontend-zvxn.vercel.app/"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	client = infra.ConnectToMongo()

	// Rate limiter configuration: 5 requests per second with a burst of 10
	limiter := middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(5),   // 5 request per second
			Burst:     10,              // Allow only 10 request at a time
			ExpiresIn: 5 * time.Minute, // Expiration for stored request data
		}),
	})

	e.Use(limiter)

	routes.SetupPostRoutes(e, client)
	routes.SetupWaifuRoutes(e, client)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "OK",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
