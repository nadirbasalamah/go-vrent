package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/nadirbasalamah/go-vrent/config"
	"github.com/nadirbasalamah/go-vrent/database"
	"github.com/nadirbasalamah/go-vrent/graph"
	"github.com/nadirbasalamah/go-vrent/graph/auth"
	"github.com/nadirbasalamah/go-vrent/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := config.Config("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.Middleware())

	err := database.Connect()
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v\n", err)
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("Connected to the database!")
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
