package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"strconv"
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
	var userDTO *CreateUserDTO
	err := ctx.BindJSON(&userDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	user, err := mapCreateUserDTOToUser(*userDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	userID, err := c.userService.Create(ctx.Request.Context(), user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"user_id": userID})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("param userID cannot be empty", nil))
		return
	}

	var user *domain.User
	ctx.BindJSON(&user)

	err := c.userService.Update(ctx.Request.Context(), *user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(204, nil)
}

func (c *UserController) GetUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("param userID cannot be empty", nil))
		return
	}

	user, err := c.userService.GetByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	userDTO := mapUserToUserDTO(*user)

	ctx.JSON(200, userDTO)
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("userID")
	if userID == "" {
		ctx.Error(domain.NewValidationError("param userID cannot be empty", nil))
		return
	}

	err := c.userService.Delete(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(204, nil)
}

func (c *UserController) SearchUser(ctx *gin.Context) {
	userFilters, err := c.parseQueryToUserFilters(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	users, err := c.userService.Search(ctx.Request.Context(), userFilters)
	if err != nil {
		ctx.Error(err)
		return
	}

	searchResult := mapSearchResultToSearchResultDTO(users, mapUsersToUserDTOs)

	ctx.JSON(200, searchResult)
}

func (c *UserController) parseQueryToUserFilters(ctx *gin.Context) (domain.UserFilters, error) {
	filters := domain.UserFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if emails := ctx.QueryArray("email"); len(emails) > 0 {
		filters.Email = emails
	}

	if firstNames := ctx.QueryArray("first_name"); len(firstNames) > 0 {
		filters.FirstName = firstNames
	}

	if userIDs := ctx.QueryArray("user_id"); len(userIDs) > 0 {
		filters.UserID = userIDs
	}

	if regions := ctx.QueryArray("region"); len(regions) > 0 {
		filters.Region = regions
	}

	if roles := ctx.QueryArray("role"); len(roles) > 0 {
		filters.Role = roles
	}

	if active := ctx.Query("active"); active != "" {
		activeBool := active == "true"
		filters.Active = &activeBool
	}

	validationErr := make([]error, 0)
	if limitParam := ctx.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("limit must be a number", nil))
		} else {
			filters.PagingFilter.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("offset must be a number", nil))
		} else {
			filters.PagingFilter.Offset = parsedOffset
		}
	}

	if len(validationErr) > 0 {
		return domain.UserFilters{}, errors.Join(validationErr...)
	}

	return filters, nil
}
