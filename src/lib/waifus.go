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
)

type Waifu struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Image       string             `bson:"image" json:"image"`
	UserId      string             `bson:"userId" json:"userId"`
	CreatedTime string             `bson:"createdTime" json:"createdTime"`
	UpdatedTime string             `bson:"updatedTime" json:"updatedTime"`
	Rating      int64              `bson:"rating" json:"rating"`
	Favorites   int64              `bson:"favorites" json:"favorites"`
	Status      int64              `bson:"status" json:"status"`
	MommyMeter  string             `bson:"mommyMeter" json:"mommyMeter"`
	Comments    string             `bson:"comments" json:"comments"`
}

func GetWaifu(c echo.Context, client *mongo.Client) error {
	collection := client.Database("animoshiApi").Collection("waifus")

	if !utils.ValidateQueryParams(c, []string{"id"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	idParam := c.QueryParam("id")
	waifuID, err := primitive.ObjectIDFromHex(idParam)

	var waifu map[string]interface{}
	err = collection.FindOne(context.TODO(), bson.M{"_id": waifuID}).Decode(&waifu)
	if err != nil {
		log.Println(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Waifu not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching waifu data"})
	}

	return c.JSON(http.StatusOK, waifu)
}

func GetWaifus(c echo.Context, client *mongo.Client) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")

	if !utils.ValidateQueryParams(c, []string{"limit", "offset"}) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid params"})
	}

	waifus, err := infra.FindAllFromCollection(infra.FindAllCollectionsParams{
		CollectionName: "waifus",
		Client:         client,
		Filter:         bson.D{{"status", "APPROVED"}},
		Limit:          limit,
		Offset:         offset,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, waifus)
}
