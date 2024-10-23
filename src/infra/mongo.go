package infra

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strconv"
	"time"
)

func ConnectToMongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

type FindAllCollectionsParams struct {
	CollectionName string
	Client         *mongo.Client
	Filter         bson.D
	Limit          string
	Offset         string
}

type CountCollectionParams struct {
	CollectionName string
	Client         *mongo.Client
	Filter         bson.D
}

func FindAllFromCollection(params FindAllCollectionsParams) ([]map[string]interface{}, error) {
	collection := params.Client.Database("animoshiApi").Collection(params.CollectionName)

	limitInt, err := strconv.Atoi(params.Limit)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nil, err
	}

	offsetInt, err := strconv.Atoi(params.Offset)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nil, err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdTime", -1}})
	findOptions.SetLimit(int64(limitInt))
	findOptions.SetSkip(int64(offsetInt))

	cur, err := collection.Find(context.TODO(), params.Filter, findOptions)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {

		}
	}(cur, context.TODO())

	var items []map[string]interface{}

	for cur.Next(context.TODO()) {
		var item map[string]interface{}
		err := cur.Decode(&item)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
		}
		items = append(items, item)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("Error converting string to int:", err)
	}

	return items, nil
}

func CountCollection(params CountCollectionParams) (int64, error) {
	collection := params.Client.Database("animoshiApi").Collection(params.CollectionName)

	count, err := collection.CountDocuments(context.TODO(), params.Filter)
	if err != nil {
		fmt.Println("Error counting collection:", err)
		return 0, err
	}

	return count, nil
}

func InsertOne(collectionName string, client *mongo.Client, item interface{}) error {
	collection := client.Database("animoshiApi").Collection(collectionName)

	_, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		fmt.Println("Error inserting item:", err)
		return err
	}

	return nil
}
