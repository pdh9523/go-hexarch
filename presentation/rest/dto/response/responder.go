package response

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/error_code"
	"net/http"
)

type Responder struct {
}

func NewResponder() *Responder {
	return &Responder{}
}

func (r *Responder) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessResponse(data))
}

func (r *Responder) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, NewSuccessResponse(data))
}

func (r *Responder) Error(c *gin.Context, err error) {
	mapping := error_code.GetErrorMapping(err)
	c.JSON(mapping.HttpStatus, NewErrorResponse(mapping.Code, err.Error()))
}
