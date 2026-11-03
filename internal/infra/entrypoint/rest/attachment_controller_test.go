package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func withRequesterID(userID string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}

func TestAttachmentController_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("deletes an attachment by id when the requester is an admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAttachmentService := mock_application.NewMockAttachmentService(ctrl)
		mockAttachmentService.EXPECT().
			GetByID(gomock.Any(), testAttachment).
			Return(&domain.Attachment{AttachmentID: testAttachment, CreatedBy: testAuthor}, nil)
		mockAttachmentService.EXPECT().DeleteByID(gomock.Any(), testAttachment, testAdmin).Return(nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testAdminID).
			Return(&domain.User{UserID: testAdminID, Username: testAdmin, Role: domain.ADMIN}, nil)

		c := NewAttachmentController(mockAttachmentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testAdminID))
		router.DELETE("/attachments/:attachmentID", c.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/attachments/attachment-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("deletes an attachment by id when the requester is its own uploader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAttachmentService := mock_application.NewMockAttachmentService(ctrl)
		mockAttachmentService.EXPECT().
			GetByID(gomock.Any(), testAttachment).
			Return(&domain.Attachment{AttachmentID: testAttachment, CreatedBy: testAuthor}, nil)
		mockAttachmentService.EXPECT().DeleteByID(gomock.Any(), testAttachment, testAuthor).Return(nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testAuthorID).
			Return(&domain.User{UserID: testAuthorID, Username: testAuthor, Role: domain.OPERATOR}, nil)

		c := NewAttachmentController(mockAttachmentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testAuthorID))
		router.DELETE("/attachments/:attachmentID", c.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/attachments/attachment-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("rejects deletion when the requester is neither an admin nor the uploader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAttachmentService := mock_application.NewMockAttachmentService(ctrl)
		mockAttachmentService.EXPECT().
			GetByID(gomock.Any(), testAttachment).
			Return(&domain.Attachment{AttachmentID: testAttachment, CreatedBy: testAuthor}, nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testOtherUserID).
			Return(&domain.User{UserID: testOtherUserID, Username: testOtherUser, Role: domain.OPERATOR}, nil)

		c := NewAttachmentController(mockAttachmentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testOtherUserID))
		router.Use(testErrorEncoder())
		router.DELETE("/attachments/:attachmentID", c.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/attachments/attachment-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
