package container

import userHandler "github.com/pdh9523/go-hexarch/presentation/rest/handler/user"

type Container interface {
	UserHandler() *userHandler.UserHandler
}

type container struct {
	userHandler *userHandler.UserHandler
}

func (c *container) UserHandler() *userHandler.UserHandler {
	return c.userHandler
}
