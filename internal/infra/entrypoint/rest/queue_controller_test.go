package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestQueueController_parseQueryToFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := QueueController{}

	tests := []struct {
		name        string
		queryParams url.Values
		wantFilters domain.QueueFilters
	}{
		{
			name:        "no filters — defaults only",
			queryParams: url.Values{},
			wantFilters: domain.QueueFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "category and state filters",
			queryParams: url.Values{"category": {"mobile"}, "state": {"SP"}},
			wantFilters: domain.QueueFilters{
				Category:     []string{"mobile"},
				State:        []string{"SP"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "active filter",
			queryParams: url.Values{"active": {"false"}},
			wantFilters: domain.QueueFilters{
				Active:       boolPtr(false),
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "custom limit and offset",
			queryParams: url.Values{"limit": {"25"}, "offset": {"5"}},
			wantFilters: domain.QueueFilters{
				PagingFilter: domain.PagingFilter{Limit: 25, Offset: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/?"+tt.queryParams.Encode(), nil)

			got := c.parseQueryToFilters(ctx)

			assert.Equal(t, tt.wantFilters.Category, got.Category)
			assert.Equal(t, tt.wantFilters.State, got.State)
			assert.Equal(t, tt.wantFilters.Active, got.Active)
			assert.Equal(t, tt.wantFilters.Limit, got.Limit)
			assert.Equal(t, tt.wantFilters.Offset, got.Offset)
		})
	}
}

func TestQueueController_SearchQueues(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		Search(gomock.Any(), domain.QueueFilters{
			Category:     []string{"mobile"},
			PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
		}).
		Return(domain.PagingResult[domain.Queue]{
			Result: []domain.Queue{{QueueID: "queue-1", Name: "SP Mobile"}},
			Paging: domain.Paging{Total: 1, Limit: 10, Offset: 0},
		}, nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.GET("/queues", c.SearchQueues)

	req := httptest.NewRequest(http.MethodGet, "/queues?category=mobile", nil)
	req = req.WithContext(context.Background())
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQueueController_CreateQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creates a queue successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockQueueService(ctrl)
		mockService.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return("queue-1", nil)

		c := NewQueueController(mockService)

		router := gin.New()
		router.POST("/queues", c.CreateQueue)

		body, _ := json.Marshal(CreateQueueDTO{
			Name:      "SP Mobile",
			Category:  "mobile",
			States:    []string{"SP"},
			CreatedBy: "author-1",
		})

		req := httptest.NewRequest(http.MethodPost, "/queues", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("does not call the service when the payload is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mock_application.NewMockQueueService(ctrl)
		c := NewQueueController(mockService)

		router := gin.New()
		router.POST("/queues", c.CreateQueue)

		body, _ := json.Marshal(CreateQueueDTO{Name: "", Category: "", CreatedBy: "author-1"})

		req := httptest.NewRequest(http.MethodPost, "/queues", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusCreated, w.Code)
	})
}

func TestQueueController_GetQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		GetByID(gomock.Any(), "queue-1").
		Return(&domain.Queue{QueueID: "queue-1", Name: "SP Mobile"}, nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.GET("/queues/:queueID", c.GetQueue)

	req := httptest.NewRequest(http.MethodGet, "/queues/queue-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQueueController_AddMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		AddMember(gomock.Any(), "queue-1", "user-1").
		Return(nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.POST("/queues/:queueID/members", c.AddMember)

	body, _ := json.Marshal(AddQueueMemberDTO{UserID: "user-1"})
	req := httptest.NewRequest(http.MethodPost, "/queues/queue-1/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestQueueController_RemoveMember(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		RemoveMember(gomock.Any(), "queue-1", "user-1").
		Return(nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.DELETE("/queues/:queueID/members/:userID", c.RemoveMember)

	req := httptest.NewRequest(http.MethodDelete, "/queues/queue-1/members/user-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestQueueController_GetMembers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		GetMembers(gomock.Any(), "queue-1").
		Return([]domain.User{{UserID: "user-1"}}, nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.GET("/queues/:queueID/members", c.GetMembers)

	req := httptest.NewRequest(http.MethodGet, "/queues/queue-1/members", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestQueueController_UpdateQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		Update(gomock.Any(), "queue-1", gomock.Any()).
		Return(nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.PUT("/queues/:queueID", c.UpdateQueue)

	body, _ := json.Marshal(UpdateQueueDTO{UpdatedBy: "author-2"})
	req := httptest.NewRequest(http.MethodPut, "/queues/queue-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestQueueController_DeleteQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_application.NewMockQueueService(ctrl)
	mockService.EXPECT().
		Delete(gomock.Any(), "queue-1").
		Return(nil)

	c := NewQueueController(mockService)

	router := gin.New()
	router.DELETE("/queues/:queueID", c.DeleteQueue)

	req := httptest.NewRequest(http.MethodDelete, "/queues/queue-1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func boolPtr(b bool) *bool {
	return &b
}
