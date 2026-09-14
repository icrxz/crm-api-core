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
)

func TestCommentController_AddAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("adds an attachment to an existing comment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockCommentService(ctrl)
		mockService.EXPECT().
			AddAttachment(gomock.Any(), testCommentID, gomock.Any()).
			Return(domain.Attachment{AttachmentID: testAttachment, CommentID: testCommentID, FileName: testFileName}, nil)

		c := NewCommentController(mockService)

		router := gin.New()
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

	t.Run("propagates the service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockCommentService(ctrl)
		mockService.EXPECT().
			AddAttachment(gomock.Any(), testCommentID, gomock.Any()).
			Return(domain.Attachment{}, domain.NewNotFoundError("comment not found", nil))

		c := NewCommentController(mockService)

		router := gin.New()
		router.Use(func(ctx *gin.Context) {
			ctx.Next()
			if len(ctx.Errors) > 0 {
				ctx.JSON(http.StatusNotFound, gin.H{"error": ctx.Errors[0].Error()})
			}
		})
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

	t.Run("updates the comment content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockCommentService(ctrl)
		mockService.EXPECT().
			UpdateContent(gomock.Any(), testCommentID, "pagamento realizado no dia 01/01/2024", testAuthor).
			Return(nil)

		c := NewCommentController(mockService)

		router := gin.New()
		router.PATCH("/comments/:commentID", c.UpdateContent)

		body, _ := json.Marshal(UpdateCommentDTO{
			Content:   "pagamento realizado no dia 01/01/2024",
			UpdatedBy: testAuthor,
		})

		req := httptest.NewRequest(http.MethodPatch, "/comments/"+testCommentID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
