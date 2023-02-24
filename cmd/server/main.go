package main

import (
	"log"
	"net/http"

	"github.com/53jk1/go-graphql-todo/internal/graphql"
	"github.com/53jk1/go-graphql-todo/internal/graphql/generated"
	"github.com/53jk1/go-graphql-todo/internal/models"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := defaultPort

	todoRepo := models.NewTodoRepository()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graphql.Resolver{TodoRepo: todoRepo}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
