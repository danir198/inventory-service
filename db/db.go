package db

import (
	"context"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewDatabase(client *mongo.Client, dbName string) *Database {
	return &Database{
		Client: client,
		DB:     client.Database(dbName),
	}
}

func (d *Database) InitializeDatabase() error {

	// Check if the database exists
	collectionNames, err := d.DB.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return err
	}

	// Create collections if they do not exist
	collectionExists := false
	for _, name := range collectionNames {
		if name == "products" {
			collectionExists = true
			break
		}
	}

	if !collectionExists {
		err := d.DB.CreateCollection(context.TODO(), "products")
		if err != nil {
			return err
		}
		log.Printf("Created  products collection db")

		// Create index on product_id field in the products collection
		collection := d.DB.Collection("products")
		_, err = collection.Indexes().CreateOne(
			context.TODO(),
			mongo.IndexModel{
				Keys:    bson.M{"product_id": 1},
				Options: options.Index().SetUnique(true),
			},
		)
		if err != nil {
			// Check if the error is because the index already exists
			if !strings.Contains(err.Error(), "already exists") {
				return fmt.Errorf("error creating index: %w", err)
			}
			log.Println("Index on product_id already exists")
		} else {
			log.Println("Index on product_id created successfully")
		}

	}

	return nil
}

func (d *Database) FindProductByID(id string) (bson.M, error) {
	var result bson.M
	collection := d.DB.Collection("products")
	err := collection.FindOne(context.TODO(), bson.M{"product_id": id}).Decode(&result)
	return result, err
}

func (d *Database) UpdateProductByID(id string, update bson.M) (*mongo.UpdateResult, error) {
	collection := d.DB.Collection("products")
	filter := bson.M{"product_id": id}
	updateDoc := bson.M{"$set": update}
	opts := options.Update().SetUpsert(true)
	return collection.UpdateOne(context.TODO(), filter, updateDoc, opts)
}

func (d *Database) CreateProduct(product bson.M) (*mongo.InsertOneResult, error) {
	collection := d.DB.Collection("products")
	return collection.InsertOne(context.TODO(), product)
}

func (d *Database) DeleteProductByID(id string) (*mongo.DeleteResult, error) {
	collection := d.DB.Collection("products")
	return collection.DeleteOne(context.TODO(), bson.M{"product_id": id})
}

func (d *Database) ListProducts() ([]bson.M, error) {
	collection := d.DB.Collection("products")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []bson.M
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (d *Database) SearchProducts(filter bson.M) ([]bson.M, error) {
	collection := d.DB.Collection("products")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []bson.M
	if err = cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}
	return products, nil
}
