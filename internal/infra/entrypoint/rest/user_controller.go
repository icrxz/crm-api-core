package rest

import "github.com/gin-gonic/gin"

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

func (c *UserController) CreateUser(ctx *gin.Context) {}

func (c *UserController) UpdateUser(ctx *gin.Context) {}

func (c *UserController) GetUser(ctx *gin.Context) {}

func (c *UserController) DeleteUser(ctx *gin.Context) {}

func (c *UserController) SearchUser(ctx *gin.Context) {}
