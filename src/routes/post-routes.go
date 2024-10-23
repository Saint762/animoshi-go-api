package routes

import (
	"animoshi-api-go/src/lib"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func SetupPostRoutes(e *echo.Echo, client *mongo.Client) {
	// GET ROUTES
	e.GET("/post", func(c echo.Context) error {
		return lib.GetPost(c, client)
	})

	e.GET("/posts", func(c echo.Context) error {
		return lib.GetPosts(c, client)
	})

	e.GET("/postsByUserId", func(c echo.Context) error {
		return lib.GetPostsByUserId(c, client)
	})

	e.GET("/postCountByUserId", func(c echo.Context) error {
		return lib.GetPostCountByUserId(c, client)
	})

	e.GET("/postComments", func(c echo.Context) error {
		return lib.GetPostCommentsByPostId(c, client)
	})

	// POST ROUTES
	e.POST("/post", func(c echo.Context) error {
		post := new(lib.Post)

		if err := c.Bind(post); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if err := lib.NewPost(c, client, post); err != nil {
			return err
		}

		return nil
	})
}
