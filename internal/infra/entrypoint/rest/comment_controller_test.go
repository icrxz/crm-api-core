package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	testCommentID  = "comment-1"
	testAttachment = "attachment-1"
	testFileName   = "comprovante.png"
	testAuthor     = "user-1"
	testAdmin      = "admin-1"
	testOtherUser  = "other-user"
)

func TestCommentController_AddAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("adds an attachment to an existing comment when the requester is its author", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, CreatedBy: testAuthor}, nil)
		mockCommentService.EXPECT().
			AddAttachment(gomock.Any(), testCommentID, gomock.Any()).
			Return(domain.Attachment{AttachmentID: testAttachment, CommentID: testCommentID, FileName: testFileName}, nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testAuthor).
			Return(&domain.User{UserID: testAuthor, Role: domain.OPERATOR}, nil)

		c := NewCommentController(mockCommentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testAuthor))
		router.POST("/comments/:commentID/attachments", c.AddAttachment)

		body, _ := json.Marshal(CreateAttachmentDTO{
			FileName:      testFileName,
			AttachmentURL: "https://s3.test/" + testFileName,
			FileExtension: "png",
			Key:           testFileName,
			CreatedBy:     testAuthor,
		})

		req := httptest.NewRequest(http.MethodPost, "/comments/"+testCommentID+"/attachments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response AttachmentDTO
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
		assert.Equal(t, testAttachment, response.AttachmentID)
	})

	t.Run("rejects adding an attachment when the requester is neither an admin nor the comment author", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, CreatedBy: testAuthor}, nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testOtherUser).
			Return(&domain.User{UserID: testOtherUser, Role: domain.OPERATOR}, nil)

		c := NewCommentController(mockCommentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testOtherUser))
		router.Use(respondWithControllerError(http.StatusForbidden))
		router.POST("/comments/:commentID/attachments", c.AddAttachment)

		body, _ := json.Marshal(CreateAttachmentDTO{FileName: testFileName})
		req := httptest.NewRequest(http.MethodPost, "/comments/"+testCommentID+"/attachments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("propagates the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(nil, domain.NewNotFoundError("comment not found", nil))

		c := NewCommentController(mockCommentService, mock_application.NewMockUserService(ctrl))

		router := gin.New()
		router.Use(respondWithControllerError(http.StatusNotFound))
		router.POST("/comments/:commentID/attachments", c.AddAttachment)

		body, _ := json.Marshal(CreateAttachmentDTO{FileName: testFileName})
		req := httptest.NewRequest(http.MethodPost, "/comments/"+testCommentID+"/attachments", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCommentController_UpdateContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("updates the comment content when the requester is an admin", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, CreatedBy: testAuthor}, nil)
		mockCommentService.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "pagamento realizado no dia 01/01/2024", testAdmin).
			Return(nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testAdmin).
			Return(&domain.User{UserID: testAdmin, Role: domain.ADMIN}, nil)

		c := NewCommentController(mockCommentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testAdmin))
		router.PATCH("/comments/:commentID", c.UpdateContent)

		body, _ := json.Marshal(UpdateCommentDTO{
			Content: "pagamento realizado no dia 01/01/2024",
		})

		req := httptest.NewRequest(http.MethodPatch, "/comments/"+testCommentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("updates the comment content when the requester is its own author", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, CreatedBy: testAuthor}, nil)
		mockCommentService.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor).
			Return(nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testAuthor).
			Return(&domain.User{UserID: testAuthor, Role: domain.OPERATOR}, nil)

		c := NewCommentController(mockCommentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testAuthor))
		router.PATCH("/comments/:commentID", c.UpdateContent)

		body, _ := json.Marshal(UpdateCommentDTO{
			Content: "pagamento realizado no dia 01/01/2024",
		})

		req := httptest.NewRequest(http.MethodPatch, "/comments/"+testCommentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("rejects updating the content when the requester is neither an admin nor the author", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCommentService := mock_application.NewMockCommentService(ctrl)
		mockCommentService.EXPECT().
			GetByID(gomock.Any(), testCommentID).
			Return(&domain.Comment{CommentID: testCommentID, CreatedBy: testAuthor}, nil)

		mockUserService := mock_application.NewMockUserService(ctrl)
		mockUserService.EXPECT().
			GetByID(gomock.Any(), testOtherUser).
			Return(&domain.User{UserID: testOtherUser, Role: domain.OPERATOR}, nil)

		c := NewCommentController(mockCommentService, mockUserService)

		router := gin.New()
		router.Use(withRequesterID(testOtherUser))
		router.Use(respondWithControllerError(http.StatusForbidden))
		router.PATCH("/comments/:commentID", c.UpdateContent)

		body, _ := json.Marshal(UpdateCommentDTO{Content: "texto malicioso"})

		req := httptest.NewRequest(http.MethodPatch, "/comments/"+testCommentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
