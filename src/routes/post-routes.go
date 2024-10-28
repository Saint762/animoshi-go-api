package routes

import (
	"animoshi-api-go/src/lib"
	"animoshi-api-go/src/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
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
		title := c.FormValue("title")
		content := c.FormValue("content")
		image := c.FormValue("image")
		nsfwToggle := c.FormValue("nsfwToggle")
		userId := c.FormValue("userId")
		recaptchaToken := c.FormValue("recaptchaToken")
		aniToken := c.FormValue("aniToken")

		if len(content) > 500 {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Content is too long"})
		}

		if len(userId) > 128 {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "UserId is too long"})
		}

		if recaptchaToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Recaptcha token is required"})
		}

		valid, err := utils.VerifyRecaptcha(recaptchaToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
		}

		if !valid {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
		}

		if aniToken == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required!"})
		}

		var fileURL string
		var imageURL string

		file, err := c.FormFile("file")
		if err != nil {
			print("No File Uploaded!")
		} else {
			if file.Size > 8*1024*1024 { // 8MB
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "File size exceeds 5MB"})
			}

			contentType := file.Header.Get("Content-Type")
			validTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/gif":  true,
				"image/webp": true,
			}

			if !validTypes[contentType] {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type. Only JPG, PNG, WEBP and GIF are allowed."})
			}

			if err == nil {
				src, err := file.Open()
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
				}
				defer src.Close()

				// Upload to S3 and get URL (assuming you have a function like uploadToS3)
				fileURL, err = utils.UploadToS3(src, file)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload file"})
				}
			}
		}

		if len(image) > 0 {
			imageURL = image
		} else if len(fileURL) > 0 {
			imageURL = fileURL
		} else {
			imageURL = ""
		}

		nsfwToggleInt, err := strconv.ParseInt(nsfwToggle, 10, 64)
		if err != nil {
			fmt.Println("Error converting nsfwToggle to int64:", err)
		}

		post := lib.PostRequest{
			Title:          title,
			Content:        content,
			Image:          imageURL,
			NsfwToggle:     nsfwToggleInt,
			UserID:         userId,
			RecaptchaToken: recaptchaToken,
			AniToken:       aniToken,
		}

		if err := lib.NewPost(c, client, &post); err != nil {
			return err
		}

		return nil
	})

	e.POST("/comment", func(c echo.Context) error {
		postComment := new(lib.PostComment)

		if err := c.Bind(postComment); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if err := lib.NewPostComment(c, client, postComment); err != nil {
			return err
		}

		return nil
	})

	e.POST("/likePost", func(c echo.Context) error {
		postLike := new(lib.PostLike)

		if err := c.Bind(postLike); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		if err := lib.LikePost(c, client, postLike); err != nil {
			return err
		}

		return nil
	})
}
