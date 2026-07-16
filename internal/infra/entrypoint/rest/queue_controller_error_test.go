package rest_test

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
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint"
	"github.com/icrxz/crm-api-core/internal/infra/entrypoint/rest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newQueueRouter(t *testing.T) (*gin.Engine, *mock_application.MockQueueService) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	mockService := mock_application.NewMockQueueService(ctrl)
	c := rest.NewQueueController(mockService)

	router := gin.New()
	router.Use(entrypoint.CustomErrorEncoder())
	router.POST("/queues", c.CreateQueue)
	router.PUT("/queues/:queueID", c.UpdateQueue)
	router.GET("/queues/:queueID", c.GetQueue)
	router.GET("/queues", c.SearchQueues)
	router.DELETE("/queues/:queueID", c.DeleteQueue)
	router.POST("/queues/:queueID/members", c.AddMember)
	router.DELETE("/queues/:queueID/members/:userID", c.RemoveMember)
	router.GET("/queues/:queueID/members", c.GetMembers)
	router.GET("/users/:userID/queues", c.GetQueuesByUser)

	return router, mockService
}

func TestQueueController_CreateQueue_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return("", errors.New("boom"))

	body, _ := json.Marshal(rest.CreateQueueDTO{Name: "SP Mobile", Category: "mobile", CreatedBy: "author-1"})
	req := httptest.NewRequest(http.MethodPost, "/queues", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_UpdateQueue_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		Update(gomock.Any(), "queue-1", gomock.Any()).
		Return(domain.NewValidationError("bad input", nil))

	body, _ := json.Marshal(rest.UpdateQueueDTO{UpdatedBy: "author-2"})
	req := httptest.NewRequest(http.MethodPut, "/queues/queue-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQueueController_UpdateQueue_InvalidBody(t *testing.T) {
	router, _ := newQueueRouter(t)

	req := httptest.NewRequest(http.MethodPut, "/queues/queue-1", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQueueController_GetQueue_NotFound(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		GetByID(gomock.Any(), "queue-1").
		Return(nil, domain.NewNotFoundError("no queue found with this id", nil))

	req := httptest.NewRequest(http.MethodGet, "/queues/queue-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestQueueController_SearchQueues_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		Search(gomock.Any(), gomock.Any()).
		Return(domain.PagingResult[domain.Queue]{}, errors.New("boom"))

	req := httptest.NewRequest(http.MethodGet, "/queues", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_DeleteQueue_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		Delete(gomock.Any(), "queue-1").
		Return(errors.New("boom"))

	req := httptest.NewRequest(http.MethodDelete, "/queues/queue-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_AddMember_InvalidBody(t *testing.T) {
	router, _ := newQueueRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/queues/queue-1/members", bytes.NewReader([]byte(`{`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestQueueController_AddMember_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		AddMember(gomock.Any(), "queue-1", "user-1").
		Return(errors.New("boom"))

	body, _ := json.Marshal(rest.AddQueueMemberDTO{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/queues/queue-1/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_RemoveMember_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		RemoveMember(gomock.Any(), "queue-1", "user-1").
		Return(errors.New("boom"))

	req := httptest.NewRequest(http.MethodDelete, "/queues/queue-1/members/user-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_GetMembers_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		GetMembers(gomock.Any(), "queue-1").
		Return(nil, errors.New("boom"))

	req := httptest.NewRequest(http.MethodGet, "/queues/queue-1/members", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestQueueController_GetQueuesByUser_ServiceError(t *testing.T) {
	router, mockService := newQueueRouter(t)

	mockService.EXPECT().
		GetQueuesByUser(gomock.Any(), "user-1").
		Return(nil, errors.New("boom"))

	req := httptest.NewRequest(http.MethodGet, "/users/user-1/queues", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
