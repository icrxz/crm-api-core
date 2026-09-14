package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAttachmentController_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("deletes an attachment by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockAttachmentService(ctrl)
		mockService.EXPECT().DeleteByID(gomock.Any(), testAttachment).Return(nil)

		c := NewAttachmentController(mockService)

		router := gin.New()
		router.DELETE("/attachments/:attachmentID", c.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/attachments/attachment-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
