package rest

import (
	"context"
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

func TestPartnerController_parseQueryToFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c := PartnerController{}

	tests := []struct {
		name        string
		queryParams url.Values
		wantFilters domain.PartnerFilters
		wantErr     bool
	}{
		{
			name:        "both name and city — sets NameOrCity with OR intent",
			queryParams: url.Values{"name": {"João"}, "city": {"São Paulo"}},
			wantFilters: domain.PartnerFilters{
				NameOrCity: &domain.NameOrCityFilter{
					Name: []string{"João"},
					City: []string{"São Paulo"},
				},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "only first_name — sets FirstName, NameOrCity nil",
			queryParams: url.Values{"first_name": {"João"}},
			wantFilters: domain.PartnerFilters{
				FirstName:    []string{"João"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "only city — sets City, NameOrCity nil",
			queryParams: url.Values{"city": {"São Paulo"}},
			wantFilters: domain.PartnerFilters{
				City:         []string{"São Paulo"},
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "neither first_name nor city — both empty",
			queryParams: url.Values{},
			wantFilters: domain.PartnerFilters{
				PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
			},
		},
		{
			name:        "custom limit and offset",
			queryParams: url.Values{"limit": {"20"}, "offset": {"5"}},
			wantFilters: domain.PartnerFilters{
				PagingFilter: domain.PagingFilter{Limit: 20, Offset: 5},
			},
		},
		{
			name:        "invalid limit returns error",
			queryParams: url.Values{"limit": {"abc"}},
			wantErr:     true,
		},
		{
			name:        "invalid offset returns error",
			queryParams: url.Values{"offset": {"xyz"}},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/?"+tt.queryParams.Encode(), nil)

			got, err := c.parseQueryToFilters(ctx)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantFilters.NameOrCity, got.NameOrCity)
			assert.Equal(t, tt.wantFilters.FirstName, got.FirstName)
			assert.Equal(t, tt.wantFilters.City, got.City)
			assert.Equal(t, tt.wantFilters.Limit, got.Limit)
			assert.Equal(t, tt.wantFilters.Offset, got.Offset)
		})
	}
}

func TestPartnerController_SearchPartners(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    url.Values
		mockSetup      func(mock *mock_application.MockPartnerService)
		wantStatusCode int
	}{
		{
			name:        "both name and city — passes NameOrCity filter to service",
			queryParams: url.Values{"name": {"João"}, "city": {"São Paulo"}},
			mockSetup: func(mock *mock_application.MockPartnerService) {
				expectedFilters := domain.PartnerFilters{
					NameOrCity: &domain.NameOrCityFilter{
						Name: []string{"João"},
						City: []string{"São Paulo"},
					},
					PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
				}
				mock.EXPECT().
					Search(gomock.Any(), expectedFilters).
					Return(domain.PagingResult[domain.Partner]{}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "only first_name — passes FirstName filter to service",
			queryParams: url.Values{"first_name": {"João"}},
			mockSetup: func(mock *mock_application.MockPartnerService) {
				expectedFilters := domain.PartnerFilters{
					FirstName:    []string{"João"},
					PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
				}
				mock.EXPECT().
					Search(gomock.Any(), expectedFilters).
					Return(domain.PagingResult[domain.Partner]{}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:        "only city — passes City filter to service",
			queryParams: url.Values{"city": {"São Paulo"}},
			mockSetup: func(mock *mock_application.MockPartnerService) {
				expectedFilters := domain.PartnerFilters{
					City:         []string{"São Paulo"},
					PagingFilter: domain.PagingFilter{Limit: 10, Offset: 0},
				}
				mock.EXPECT().
					Search(gomock.Any(), expectedFilters).
					Return(domain.PagingResult[domain.Partner]{}, nil)
			},
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_application.NewMockPartnerService(ctrl)
			tt.mockSetup(mockService)

			c := NewPartnerController(mockService)

			router := gin.New()
			router.GET("/partners", c.SearchPartners)

			req := httptest.NewRequest(http.MethodGet, "/partners?"+tt.queryParams.Encode(), nil)
			req = req.WithContext(context.Background())
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}
