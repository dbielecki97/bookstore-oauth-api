package http

import (
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	GetById(*gin.Context)
}

type handler struct {
	service token.Service
	router  *gin.Engine
}

func NewHandler(service token.Service) Handler {
	return &handler{service: service}
}

func (h handler) GetById(c *gin.Context) {
	tid := c.Param("token_id")
	t, err := h.service.GetTokenById(tid)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, t)
}
