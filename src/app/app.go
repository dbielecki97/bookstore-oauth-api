package app

import (
	client "github.com/dbielecki97/bookstore-oauth-api/src/clients/cassandra"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-oauth-api/src/http"
	repo "github.com/dbielecki97/bookstore-oauth-api/src/repository/cassandra"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	s := token.NewService(repo.New(client.GetSession()))
	h := http.NewHandler(s)

	router.GET("/oauth/token/:token_id", h.GetById)
	router.POST("/oauth/token", h.Create)

	router.Run(":8081")
}
