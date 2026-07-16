package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func withUserID(userID string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

// testErrorEncoder mirrors entrypoint.CustomErrorEncoder to translate ctx.Error
// into the HTTP status carried by the domain.CustomError, without importing the
// parent entrypoint package (which itself imports rest, causing an import cycle).
func testErrorEncoder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}

		var customErr *domain.CustomError
		if errors.As(ctx.Errors[0].Err, &customErr) {
			ctx.AbortWithStatusJSON(customErr.StatusCode(), gin.H{"message": customErr.Error()})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "unexpected error"})
	}
}

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creates a user successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return("user-1", nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.POST("/users", c.CreateUser)

		body, _ := json.Marshal(CreateUserDTO{
			Username:  "johndoe",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Role:      domain.OPERATOR,
			Region:    1,
			Password:  "s3cr3t",
			CreatedBy: "author-1",
		})

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("returns error when body is malformed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.POST("/users", c.CreateUser)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return("", domain.NewConflictError("email already in use", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.POST("/users", c.CreateUser)

		body, _ := json.Marshal(CreateUserDTO{
			Username:  "johndoe",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Role:      domain.OPERATOR,
			Region:    1,
			Password:  "s3cr3t",
			CreatedBy: "author-1",
		})

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("admin updates another user successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "admin-1").
			Return(&domain.User{UserID: "admin-1", Role: domain.ADMIN}, nil)
		mockService.EXPECT().
			Update(gomock.Any(), "user-1", "author-2", gomock.Any()).
			Return(nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.PUT("/users/:userID", withUserID("admin-1"), c.UpdateUser)

		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "author-2"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("operator updates their own profile fields successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-1").
			Return(&domain.User{UserID: "user-1", Role: domain.OPERATOR}, nil)
		mockService.EXPECT().
			Update(gomock.Any(), "user-1", "user-1", gomock.Any()).
			Return(nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.PUT("/users/:userID", withUserID("user-1"), c.UpdateUser)

		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "user-1"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("rejects when an operator tries to update another user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-2").
			Return(&domain.User{UserID: "user-2", Role: domain.OPERATOR}, nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID", withUserID("user-2"), c.UpdateUser)

		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "user-2"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("rejects when an operator tries to change role, region or active status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-1").
			Return(&domain.User{UserID: "user-1", Role: domain.OPERATOR}, nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID", withUserID("user-1"), c.UpdateUser)

		adminRole := domain.ADMIN
		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "user-1", Role: &adminRole})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns error when the requester cannot be found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-1").
			Return(nil, domain.NewNotFoundError("user not found", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID", withUserID("user-1"), c.UpdateUser)

		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "user-1"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "admin-1").
			Return(&domain.User{UserID: "admin-1", Role: domain.ADMIN}, nil)
		mockService.EXPECT().
			Update(gomock.Any(), "user-1", "author-2", gomock.Any()).
			Return(domain.NewNotFoundError("user not found", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID", withUserID("admin-1"), c.UpdateUser)

		body, _ := json.Marshal(UpdateUserDTO{UpdatedBy: "author-2"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserController_GetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-1").
			Return(&domain.User{UserID: "user-1", Username: "johndoe"}, nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.GET("/users/:userID", c.GetUser)

		req := httptest.NewRequest(http.MethodGet, "/users/user-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("returns error when user is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			GetByID(gomock.Any(), "user-1").
			Return(nil, domain.NewNotFoundError("user not found", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.GET("/users/:userID", c.GetUser)

		req := httptest.NewRequest(http.MethodGet, "/users/user-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns error when userID param is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.GET("/users", c.GetUser)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("deletes the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Delete(gomock.Any(), "user-1").
			Return(nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.DELETE("/users/:userID", c.DeleteUser)

		req := httptest.NewRequest(http.MethodDelete, "/users/user-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Delete(gomock.Any(), "user-1").
			Return(domain.NewNotFoundError("user not found", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.DELETE("/users/:userID", c.DeleteUser)

		req := httptest.NewRequest(http.MethodDelete, "/users/user-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns error when userID param is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.DELETE("/users", c.DeleteUser)

		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_SearchUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns matching users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Search(gomock.Any(), domain.UserFilters{
				Email:        []string{"john@doe.com"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			}).
			Return(domain.PagingResult[domain.User]{
				Result: []domain.User{{UserID: "user-1", Username: "johndoe"}},
				Paging: domain.Paging{Total: 1, Limit: 10, Offset: 0},
			}, nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.GET("/users", c.SearchUser)

		req := httptest.NewRequest(http.MethodGet, "/users?email=john@doe.com", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("returns error when limit is not a number", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.GET("/users", c.SearchUser)

		req := httptest.NewRequest(http.MethodGet, "/users?limit=abc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns error when offset is not a number", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.GET("/users", c.SearchUser)

		req := httptest.NewRequest(http.MethodGet, "/users?offset=abc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Search(gomock.Any(), gomock.Any()).
			Return(domain.PagingResult[domain.User]{}, domain.NewValidationError("invalid filters", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.GET("/users", c.SearchUser)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("applies first_name, user_id, region, role and active filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			Search(gomock.Any(), domain.UserFilters{
				FirstName:    []string{"John"},
				UserID:       []string{"user-1"},
				Region:       []string{"1"},
				Role:         []string{"operator"},
				Active:       boolPtr(true),
				PagingFilter: domain.PagingFilter{Limit: 25, Offset: 5},
			}).
			Return(domain.PagingResult[domain.User]{}, nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.GET("/users", c.SearchUser)

		req := httptest.NewRequest(http.MethodGet, "/users?first_name=John&user_id=user-1&region=1&role=operator&active=true&limit=25&offset=5", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUserController_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("changes the password when the requester is the owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			ChangePassword(gomock.Any(), "user-1", "old-password", "new-password").
			Return(nil)

		c := NewUserController(mockService)

		router := gin.New()
		router.PUT("/users/:userID/password", withUserID("user-1"), c.ChangePassword)

		body, _ := json.Marshal(ChangePasswordDTO{OldPassword: "old-password", NewPassword: "new-password"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("rejects when the requester is not the owner", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID/password", withUserID("user-2"), c.ChangePassword)

		body, _ := json.Marshal(ChangePasswordDTO{OldPassword: "old-password", NewPassword: "new-password"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("returns error when payload is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)

		c := NewUserController(mockService)

		router := gin.New()
		router.PUT("/users/:userID/password", withUserID("user-1"), c.ChangePassword)

		body, _ := json.Marshal(ChangePasswordDTO{OldPassword: "", NewPassword: ""})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns error when userID param is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/password", withUserID("user-1"), c.ChangePassword)

		req := httptest.NewRequest(http.MethodPut, "/password", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockUserService(ctrl)
		mockService.EXPECT().
			ChangePassword(gomock.Any(), "user-1", "old-password", "new-password").
			Return(domain.NewValidationError("old_password is incorrect", nil))

		c := NewUserController(mockService)

		router := gin.New()
		router.Use(testErrorEncoder())
		router.PUT("/users/:userID/password", withUserID("user-1"), c.ChangePassword)

		body, _ := json.Marshal(ChangePasswordDTO{OldPassword: "old-password", NewPassword: "new-password"})
		req := httptest.NewRequest(http.MethodPut, "/users/user-1/password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
