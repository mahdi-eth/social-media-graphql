package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	graph "github.com/mahdi-eth/social-media-graphql/api/graphql"
	"github.com/mahdi-eth/social-media-graphql/internal/db"
)


func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	
	port := os.Getenv("PORT")

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(&transport.Websocket{})

	db.Connect()
	fmt.Println("Successfully connected to db")

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
