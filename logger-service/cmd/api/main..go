package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	gRpcPort = "50001"
	mongoUrl = "mongodb://mongo:27017"
)

var mongoClient *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	client, err := connectToMongoDB()

	if err != nil {
		log.Panic(err)
	}

	mongoClient = client

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		//close connection
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Panic(err)
		}
		log.Println("Connection to MongoDB closed.")
	}()

	app := Config{
		Models: data.New(mongoClient),
	}

	// app.serve()
	log.Printf("Running logger service on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		println(err)
		log.Panic(err)
	}

}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}

// 	err := srv.ListenAndServe()

// 	if err != nil {
// 		println(err)
// 		log.Panic(err)
// 	}

// 	log.Printf("Running logger service on port %s\n", webPort)
// }

func connectToMongoDB() (*mongo.Client, error) {

	mongoClientOptions := options.Client().ApplyURI(mongoUrl)
	//TODO get from env
	mongoClientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	client, err := mongo.Connect(context.TODO(), mongoClientOptions)

	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")

	return client, nil
}
