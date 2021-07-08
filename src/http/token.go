package http

import (
	"errors"
	"github.com/dbielecki97/bookstore-oauth-api/src/domain/token"
	"github.com/dbielecki97/bookstore-utils-go/errs"
	"github.com/dbielecki97/bookstore-utils-go/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler interface {
	GetById(*gin.Context)
	Create(*gin.Context)
	UpdateExpiration(*gin.Context)
}

type handler struct {
	service token.Service
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

func (h handler) Create(c *gin.Context) {
	var r token.Request
	if err := c.ShouldBindJSON(&r); err != nil {
		restErr := errs.NewBadRequestErr("invalid json body")
		c.JSON(restErr.StatusCode, restErr)
		return
	}

	t, err := h.service.CreateToken(r)
	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusCreated, t)
}

func (h handler) UpdateExpiration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, map[string]interface{}{"info": "not implemented"})
	logger.Error("used not implemented endpoint", errors.New(c.Request.RequestURI))
}
