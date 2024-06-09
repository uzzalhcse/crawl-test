package ninjacrawler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

type client struct {
	*mongo.Client
}

func connectDB() *client {
	return &client{
		mustGetClient(),
	}
}

func mustGetClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseURL := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		app.Config.Env("DB_USERNAME"),
		app.Config.Env("DB_PASSWORD"),
		app.Config.Env("DB_HOST"),
		app.Config.Env("DB_PORT"),
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURL))
	if err != nil {
		panic(err)
	}
	// Check if the connection is established
	err = client.Ping(ctx, nil)
	if err != nil {
		app.Logger.Error("Failed to ping MongoDB: %v", err)
	}

	return client
}
func (client *client) getCollection(collectionName string) *mongo.Collection {
	collection := client.Database(app.Name).Collection(collectionName)
	ensureUniqueIndex(collection)
	return collection
}

func ensureUniqueIndex(collection *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"url": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create index: %v", err))
	}
}
func (client *client) insert(urlCollections []UrlCollection, parent string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var documents []interface{}
	for _, urlCollection := range urlCollections {
		urlCollection := UrlCollection{
			Url:       urlCollection.Url,
			Parent:    parent,
			Status:    false,
			Error:     false,
			MetaData:  urlCollection.MetaData,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		}
		documents = append(documents, urlCollection)
	}

	opts := options.InsertMany().SetOrdered(false)

	collection := client.getCollection(app.GetCollection())
	_, _ = collection.InsertMany(ctx, documents, opts)

}
func (client *client) newSite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var documents interface{}

	documents = SiteCollection{
		Url:       app.Url,
		BaseUrl:   app.BaseUrl,
		Status:    false,
		Attempts:  0,
		StartedAt: time.Now(),
		EndedAt:   nil,
	}

	collection := client.getCollection(app.GetCollection())
	_, _ = collection.InsertOne(ctx, documents)

}
func (client *client) saveProductDetail(productDetail *ProductDetail) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := client.getCollection(app.GetCollection())
	collection.ReplaceOne(ctx, bson.D{{Key: "url", Value: productDetail.Url}}, productDetail, options.Replace().SetUpsert(true))
	defer cancel()

}

func (client *client) getUrlsFromCollection(collection string) []string {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 5}}},
	}
	return extractUrls(filterData(filterCondition, client.getCollection(collection)))
}

func (client *client) getUrlCollections(collection string) []UrlCollection {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 5}}},
	}
	return filterUrlData(filterCondition, client.getCollection(collection))
}

func filterUrlData(filterCondition bson.D, mongoCollection *mongo.Collection) []UrlCollection {
	findOptions := options.Find().SetLimit(1000) // TODO: need to refactor

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	var results []UrlCollection
	if err = cursor.All(context.TODO(), &results); err != nil {
		slog.Error(err.Error())
	}

	return results
}
func filterData(filterCondition bson.D, mongoCollection *mongo.Collection) []bson.M {
	findOptions := options.Find().SetLimit(1000) // TODO: need to refactor

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		slog.Error(err.Error())
	}

	return results
}
func extractUrls(results []bson.M) []string {
	var urls []string
	for _, result := range results {
		if url, ok := result["url"].(string); ok {
			urls = append(urls, url)
		}
	}
	return urls
}
func (client *client) close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}
