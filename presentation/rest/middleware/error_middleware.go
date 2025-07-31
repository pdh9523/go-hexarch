package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	"log"
)

func RecoveryMiddleware(responder *response.Responder) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		var err error

		switch v := recovered.(type) {
		case error:
			err = v
		case string:
			err = errors.New(v)
		default:
			err = fmt.Errorf("panic recovered: %v", v)
		}

		log.Printf("Panic recovered: %v", err)
		responder.Error(c, err)
	})
}
