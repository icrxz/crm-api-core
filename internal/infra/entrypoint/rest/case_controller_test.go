package rest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCaseController_parseQueryToFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		query    url.Values
		expected domain.CaseFilters
	}{
		{
			name:  "single case_id",
			query: url.Values{"case_id": {"abc-123"}},
			expected: domain.CaseFilters{
				CaseID:       []string{"abc-123"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "multiple case_ids comma-separated",
			query: url.Values{"case_id": {"abc-123,def-456,ghi-789"}},
			expected: domain.CaseFilters{
				CaseID:       []string{"abc-123", "def-456", "ghi-789"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "case_id combined with other filters",
			query: url.Values{"case_id": {"abc-123,def-456"}, "status": {"New"}, "region": {"1"}},
			expected: domain.CaseFilters{
				CaseID:       []string{"abc-123", "def-456"},
				Status:       []string{"New"},
				Region:       []string{"1"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "no case_id param",
			query: url.Values{"status": {"New"}},
			expected: domain.CaseFilters{
				Status:       []string{"New"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "empty query returns defaults",
			query: url.Values{},
			expected: domain.CaseFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "custom paging",
			query: url.Values{"case_id": {"abc-123"}, "limit": {"20"}, "offset": {"5"}},
			expected: domain.CaseFilters{
				CaseID:       []string{"abc-123"},
				PagingFilter: domain.PagingFilter{Limit: 20, Offset: 5, SortBy: "created_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "sort_by updated_at",
			query: url.Values{"sort_by": {"updated_at"}},
			expected: domain.CaseFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "updated_at", SortOrder: "DESC"},
			},
		},
		{
			name:  "sort_by updated_at with sort_order ASC",
			query: url.Values{"sort_by": {"updated_at"}, "sort_order": {"asc"}},
			expected: domain.CaseFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "updated_at", SortOrder: "ASC"},
			},
		},
		{
			name:  "invalid sort_by falls back to default",
			query: url.Values{"sort_by": {"drop table cases"}},
			expected: domain.CaseFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "drop table cases", SortOrder: "DESC"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/?"+tc.query.Encode(), nil)

			controller := CaseController{}
			filters := controller.parseQueryToFilters(ctx)

			assert.Equal(t, tc.expected, filters)
		})
	}
}

func TestCaseController_SearchCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		query           url.Values
		expectedFilters domain.CaseFilters
		serviceResult   domain.PagingResult[domain.Case]
		expectedStatus  int
	}{
		{
			name:  "filters by comma-separated case_ids",
			query: url.Values{"case_id": {"id-1,id-2"}},
			expectedFilters: domain.CaseFilters{
				CaseID:       []string{"id-1", "id-2"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
			serviceResult:  domain.PagingResult[domain.Case]{Result: []domain.Case{}, Paging: domain.Paging{Total: 0, Limit: 10, Offset: 0}},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "filters by single case_id",
			query: url.Values{"case_id": {"id-1"}},
			expectedFilters: domain.CaseFilters{
				CaseID:       []string{"id-1"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
			serviceResult:  domain.PagingResult[domain.Case]{Result: []domain.Case{}, Paging: domain.Paging{Total: 0, Limit: 10, Offset: 0}},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "no case_id filter passes empty slice",
			query: url.Values{},
			expectedFilters: domain.CaseFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0, SortBy: "created_at", SortOrder: "DESC"},
			},
			serviceResult:  domain.PagingResult[domain.Case]{Result: []domain.Case{}, Paging: domain.Paging{Total: 0, Limit: 10, Offset: 0}},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCaseService := mock_application.NewMockCaseService(ctrl)
			mockCaseService.EXPECT().
				SearchCases(gomock.Any(), tc.expectedFilters).
				Return(tc.serviceResult, nil)

			w := httptest.NewRecorder()
			ctx, engine := gin.CreateTestContext(w)

			engine.GET("/cases", func(c *gin.Context) {
				controller := CaseController{caseService: mockCaseService}
				controller.SearchCases(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/cases?"+tc.query.Encode(), nil)
			ctx.Request = req
			engine.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
