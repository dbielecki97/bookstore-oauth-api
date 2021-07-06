package app

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-oauth-api/src/http"
	"github.com/dbielecki97/bookstore-oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	s := token.NewService(db.NewRepo())
	h := http.NewHandler(s)

	router.GET("/oauth/token/:token_id", h.GetById)

	router.Run(":8081")
}
