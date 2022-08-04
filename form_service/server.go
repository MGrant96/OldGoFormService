package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/polyloop/formservice/graph"
	"github.com/polyloop/formservice/graph/generated"
)


func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	defaultPort := ":8080"

	r := gin.New()
	r.Use(cors.Middleware(cors.Config{
		Origins:        "*",
		Methods:        "*",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		ValidateHeaders: false,
	}))
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run(defaultPort)
}
