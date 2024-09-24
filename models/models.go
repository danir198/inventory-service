// Step 2: Define main types and configuration
// models/models.go
package models

type Product struct {
	ID       string  `bson:"product_id" json:"id"`
	Name     string  `bson:"name" json:"name"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
}

type Config struct {
	MongoURI      string `envconfig:"MONGO_URI" required:"true"`
	DatabaseName  string `envconfig:"DB_NAME" default:"inventory"`
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:":8080"`
}
