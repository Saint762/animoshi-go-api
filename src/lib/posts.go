package lib

import (
	"animoshi-api-go/src/infra"
	"animoshi-api-go/src/utils"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Post struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title          string             `bson:"title" json:"title"`
	Content        string             `bson:"content" json:"content"`
	Image          string             `bson:"image" json:"image"`
	Video          string             `bson:"video" json:"video"`
	Likes          int64              `bson:"likes" json:"likes"`
	NsfwToggle     int64              `bson:"nsfwToggle" json:"nsfwToggle"`
	Comments       int64              `bson:"comments" json:"comments"`
	UserID         string             `bson:"userId" json:"userId"`
	UserName       string             `bson:"userName" json:"userName"`
	CreatedTime    string             `bson:"createdTime" json:"createdTime"`
	UpdatedTime    string             `bson:"updatedTime" json:"updatedTime"`
	RecaptchaToken string             `bson:"recaptchaToken,omitempty" json:"recaptchaToken"`
}

func GetPost(c echo.Context, client *mongo.Client) error {
	collection := client.Database("animoshiApi").Collection("posts")

	if !utils.ValidateQueryParams(c, []string{"id"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	idParam := c.QueryParam("id")
	postID, err := primitive.ObjectIDFromHex(idParam)

	var post map[string]interface{}
	err = collection.FindOne(context.TODO(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		log.Println(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching post data"})
	}

	return c.JSON(http.StatusOK, post)
}

func GetPosts(c echo.Context, client *mongo.Client) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")

	if !utils.ValidateQueryParams(c, []string{"limit", "offset"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	posts, err := infra.FindAllFromCollection(infra.FindAllCollectionsParams{
		CollectionName: "posts",
		Client:         client,
		Filter:         bson.D{},
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func GetPostsByUserId(c echo.Context, client *mongo.Client) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	userId := c.QueryParam("userId")

	if !utils.ValidateQueryParams(c, []string{"limit", "offset", "userId"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	posts, err := infra.FindAllFromCollection(infra.FindAllCollectionsParams{
		CollectionName: "posts",
		Client:         client,
		Filter:         bson.D{{"userId", userId}},
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, posts)
}

func GetPostCountByUserId(c echo.Context, client *mongo.Client) error {
	userId := c.QueryParam("userId")

	if !utils.ValidateQueryParams(c, []string{"userId"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	postCount, err := infra.CountCollection(infra.CountCollectionParams{
		CollectionName: "posts",
		Client:         client,
		Filter:         bson.D{{"userId", userId}},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, postCount)
}

func GetPostCommentsByPostId(c echo.Context, client *mongo.Client) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	postId := c.QueryParam("postId")

	if !utils.ValidateQueryParams(c, []string{"postId", "limit", "offset"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	postComments, err := infra.FindAllFromCollection(infra.FindAllCollectionsParams{
		CollectionName: "postComments",
		Client:         client,
		Filter:         bson.D{{"postId", postId}},
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, postComments)
}

func NewPost(c echo.Context, client *mongo.Client, post *Post) error {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	if post.RecaptchaToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Recaptcha token is required"})
	}

	valid, err := utils.VerifyRecaptcha(post.RecaptchaToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	if !valid {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	post.CreatedTime = strconv.FormatInt(currentTime, 10)
	post.UpdatedTime = strconv.FormatInt(currentTime, 10)

	post.Likes = 0
	post.Comments = 0
	post.NsfwToggle = 0
	post.UserName = post.UserID

	post.ID = primitive.NewObjectID()

	insertErr := infra.InsertOne("posts", client, post)
	if insertErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertErr.Error()})
	}

	return c.JSON(http.StatusOK, post)
}
