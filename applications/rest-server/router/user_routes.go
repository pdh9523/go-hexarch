package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pdh9523/go-hexarch/presentation/rest/dto/response"
	handler "github.com/pdh9523/go-hexarch/presentation/rest/handler/user"
	"github.com/pdh9523/go-hexarch/presentation/rest/middleware"
	"github.com/pdh9523/go-hexarch/shared/security"
)

func setupUserRoutesInternal(rg *gin.RouterGroup, userHandler *handler.UserHandler, responder *response.Responder, securityManager security.Manager) {
	users := rg.Group("/users")
	{
		users.POST("/signup", userHandler.SignUp)
		users.POST("/signin", userHandler.SignIn)
		users.GET("/username/check", userHandler.CheckUsernameAvailability)
		users.GET("/nickname/check", userHandler.CheckNicknameAvailability)

		authenticated := users.Group("")
		authenticated.Use(middleware.AuthMiddleware(responder, securityManager))
		{
			authenticated.GET("/me", userHandler.GetMyInfo)
			authenticated.POST("/signout", userHandler.SignOut)
			authenticated.PUT("/password", userHandler.ChangePassword)
			authenticated.PUT("/nickname", userHandler.ChangeNickname)
		}
	}
}
