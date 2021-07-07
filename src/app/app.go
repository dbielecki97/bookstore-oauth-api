package app

import (
	client "github.com/dbielecki97/bookstore-oauth-api/src/clients/cassandra"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-oauth-api/src/http"
	tokenRepo "github.com/dbielecki97/bookstore-oauth-api/src/repository/cassandra"
	restRepo "github.com/dbielecki97/bookstore-oauth-api/src/repository/rest"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	s := token.NewService(tokenRepo.New(client.GetSession()), restRepo.New())
	h := http.NewHandler(s)

	router.GET("/oauth/token/:token_id", h.GetById)
	router.POST("/oauth/token", h.Create)

	router.Run(":8081")
}
