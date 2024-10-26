package main

import (
	"animoshi-api-go/src/infra"
	"animoshi-api-go/src/routes"
	"context"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"net/url"
	"time"
)

var client *mongo.Client

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

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
			Rate:      rate.Limit(1),
			Burst:     10,
			ExpiresIn: 5 * time.Minute,
		}),
	})

	issuerURL, err := url.Parse("https://" + "casualhoodlums.auth0.com" + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	_, err = validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{"https://ch-auth/api"},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}

	//_, err = jwtValidator.ValidateToken(context.TODO(), "")
	//if err != nil {
	//	log.Fatalf("failed to validate token: %v", err)
	//}

	e.Use(limiter)

	routes.SetupPostRoutes(e, client)
	routes.SetupWaifuRoutes(e, client)

	e.GET("/", func(c echo.Context) error {
		// Return a JSON response with a status message
		return c.JSON(http.StatusOK, map[string]string{
			"message": "OK",
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
