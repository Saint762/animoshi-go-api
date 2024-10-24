package lib

import (
	"animoshi-api-go/src/infra"
	"animoshi-api-go/src/utils"
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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

type PostComment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PostId         string             `bson:"postId" json:"postId"`
	UserID         string             `bson:"userId" json:"userId"`
	Text           string             `bson:"text" json:"text"`
	CreatedTime    string             `bson:"createdTime" json:"createdTime"`
	UpdatedTime    string             `bson:"updatedTime" json:"updatedTime"`
	RecaptchaToken string             `bson:"recaptchaToken,omitempty" json:"recaptchaToken"`
}

type PostLike struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PostId         string             `bson:"postId" json:"postId"`
	UserID         string             `bson:"userId" json:"userId"`
	UserIP         string             `bson:"userIp" json:"userIp"`
	CreatedTime    string             `bson:"createdTime" json:"createdTime"`
	UpdatedTime    string             `bson:"updatedTime" json:"updatedTime"`
	RecaptchaToken string             `bson:"recaptchaToken,omitempty" json:"recaptchaToken"`
}

func sanitizeInput(input string) string {
	re := regexp.MustCompile(`[<>]`)
	return re.ReplaceAllString(input, "")
}

var validate = validator.New()

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

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	if limitInt > 20 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Limit cant be more than 20"})
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

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	if limitInt > 20 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Limit cant be more than 20"})
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

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	if limitInt > 20 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Limit cant be more than 20"})
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

	if err := validate.Struct(post); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	if len(post.Title) > 100 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Title is too long"})
	}

	if len(post.Content) > 500 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Content is too long"})
	}

	if len(post.Image) > 500 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Image is too long"})
	}

	if len(post.Image) > 0 && !strings.HasPrefix(post.Image, "https://") {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Image URL must start with https://"})
	}

	if len(post.UserID) > 128 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "UserId is too long"})
	}

	post.Title = sanitizeInput(post.Title)
	post.Content = sanitizeInput(post.Content)
	post.UserID = sanitizeInput(post.UserID)

	post.CreatedTime = strconv.FormatInt(currentTime, 10)
	post.UpdatedTime = strconv.FormatInt(currentTime, 10)

	post.Likes = 0
	post.Comments = 0
	post.UserName = post.UserID

	post.ID = primitive.NewObjectID()

	insertErr := infra.InsertOne("posts", client, post)
	if insertErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertErr.Error()})
	}

	return c.JSON(http.StatusOK, post)
}

func NewPostComment(c echo.Context, client *mongo.Client, postComment *PostComment) error {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	if postComment.RecaptchaToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Recaptcha token is required"})
	}

	valid, err := utils.VerifyRecaptcha(postComment.RecaptchaToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	if !valid {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	if err := validate.Struct(postComment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	if len(postComment.Text) > 500 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Content is too long"})
	}

	if len(postComment.UserID) > 128 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "UserId is too long"})
	}

	postComment.CreatedTime = strconv.FormatInt(currentTime, 10)
	postComment.UpdatedTime = strconv.FormatInt(currentTime, 10)

	postComment.ID = primitive.NewObjectID()

	collection := client.Database("animoshiApi").Collection("posts")

	postObjectID, err := primitive.ObjectIDFromHex(postComment.PostId)
	if err != nil {
		return err
	}

	var post bson.M
	err = collection.FindOne(context.TODO(), bson.M{"_id": postObjectID}).Decode(&post)
	if err != nil {
		print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Post not found"})
	}

	currentCommentCount, ok := post["comments"].(int64)
	if !ok {
		currentCommentCount = 0
	}

	newCommentCount := currentCommentCount + 1

	update := bson.M{
		"$set": bson.M{
			"comments":    newCommentCount,
			"updatedTime": strconv.FormatInt(currentTime, 10),
		},
	}

	insertErr := infra.InsertOne("postComments", client, postComment)
	if insertErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertErr.Error()})
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": postObjectID}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update comment count"})
	}

	return c.JSON(http.StatusOK, postComment)
}

func LikePost(c echo.Context, client *mongo.Client, postLike *PostLike) error {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)

	if postLike.RecaptchaToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Recaptcha token is required"})
	}

	valid, err := utils.VerifyRecaptcha(postLike.RecaptchaToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	if !valid {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid Recaptcha Token"})
	}

	if err := validate.Struct(postLike); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	postLike.UserID = "Anonymous"
	postLike.UserIP = utils.GetUserIP(c)
	postLike.CreatedTime = strconv.FormatInt(currentTime, 10)
	postLike.UpdatedTime = strconv.FormatInt(currentTime, 10)

	postLike.ID = primitive.NewObjectID()

	collection := client.Database("animoshiApi").Collection("posts")

	postObjectID, err := primitive.ObjectIDFromHex(postLike.PostId)
	if err != nil {
		return err
	}

	var post bson.M
	err = collection.FindOne(context.TODO(), bson.M{"_id": postObjectID}).Decode(&post)
	if err != nil {
		print(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Post not found"})
	}

	currentLikeCount, ok := post["likes"].(int64)
	if !ok {
		currentLikeCount = 0
	}

	newLikeCount := currentLikeCount + 1

	update := bson.M{
		"$set": bson.M{
			"likes":       newLikeCount,
			"updatedTime": strconv.FormatInt(currentTime, 10),
		},
	}

	insertErr := infra.InsertOne("postLikes", client, postLike)
	if insertErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertErr.Error()})
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": postObjectID}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update comment count"})
	}

	return c.JSON(http.StatusOK, postLike)
}
