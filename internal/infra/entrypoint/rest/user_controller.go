package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type UserController struct {
	userService application.UserService
}

func NewUserController(userService application.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var user *domain.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.Error(err)
		return
	}

	userID, err := c.userService.Create(ctx, *user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"user_id": userID})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var user *domain.User
	ctx.BindJSON(&user)

	err := c.userService.Update(ctx, *user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(204, nil)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	user, err := c.userService.GetByID(ctx, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, user)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	err := c.userService.Delete(ctx, userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(204, nil)
}

func (c *UserController) SearchUser(ctx *gin.Context) {
	var filters domain.UserFilters
	err := ctx.BindQuery(&filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	users, err := c.userService.Search(ctx, filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, users)
}
