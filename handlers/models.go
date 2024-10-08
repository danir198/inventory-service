// Step 2: Define main types and configuration
// models/models.go
package handlers

type Product struct {
	ProductID string  `bson:"product_id" json:"product_id"`
	Name      string  `bson:"name" json:"name"`
	Quantity  int     `bson:"quantity" json:"quantity"`
	Price     float64 `bson:"price" json:"price"`
}

type Config struct {
	MongoURI      string `envconfig:"MONGO_URI" required:"true"`
	DatabaseName  string `envconfig:"DB_NAME" default:"inventory"`
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:":8080"`
}
