package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildLikePart(t *testing.T) {
	tests := []struct {
		name     string
		filters  []string
		startIdx int
		key      string
		want     string
	}{
		{
			name:     "single value",
			filters:  []string{"João"},
			startIdx: 1,
			key:      "first_name",
			want:     "first_name ILIKE '%' || $1::text || '%'",
		},
		{
			name:     "multiple values",
			filters:  []string{"João", "Maria"},
			startIdx: 3,
			key:      "shipping_city",
			want:     "shipping_city ILIKE '%' || $3::text$4::text || '%'",
		},
		{
			name:     "startIdx accounts for existing args",
			filters:  []string{"SP"},
			startIdx: 5,
			key:      "shipping_city",
			want:     "shipping_city ILIKE '%' || $5::text || '%'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildLikePart(tt.filters, tt.startIdx, tt.key)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPrepareOrLikeQuery(t *testing.T) {
	tests := []struct {
		name      string
		filters1  []string
		filters2  []string
		initQuery []string
		initArgs  []any
		key1      string
		key2      string
		wantQuery []string
		wantArgs  []any
	}{
		{
			name:      "both filters set — generates OR clause",
			filters1:  []string{"João"},
			filters2:  []string{"São Paulo"},
			initQuery: []string{"1=1"},
			initArgs:  []any{},
			key1:      "first_name",
			key2:      "shipping_city",
			wantQuery: []string{
				"1=1",
				"(first_name ILIKE '%' || $1::text || '%' OR shipping_city ILIKE '%' || $2::text || '%')",
			},
			wantArgs: []any{"João", "São Paulo"},
		},
		{
			name:      "first filter empty — no-op",
			filters1:  []string{},
			filters2:  []string{"São Paulo"},
			initQuery: []string{"1=1"},
			initArgs:  []any{},
			key1:      "first_name",
			key2:      "shipping_city",
			wantQuery: []string{"1=1"},
			wantArgs:  []any{},
		},
		{
			name:      "second filter empty — no-op",
			filters1:  []string{"João"},
			filters2:  []string{},
			initQuery: []string{"1=1"},
			initArgs:  []any{},
			key1:      "first_name",
			key2:      "shipping_city",
			wantQuery: []string{"1=1"},
			wantArgs:  []any{},
		},
		{
			name:      "both filters empty — no-op",
			filters1:  []string{},
			filters2:  []string{},
			initQuery: []string{"1=1"},
			initArgs:  []any{},
			key1:      "first_name",
			key2:      "shipping_city",
			wantQuery: []string{"1=1"},
			wantArgs:  []any{},
		},
		{
			name:      "arg indices account for pre-existing args",
			filters1:  []string{"João"},
			filters2:  []string{"São Paulo"},
			initQuery: []string{"1=1", "active = $1"},
			initArgs:  []any{true},
			key1:      "first_name",
			key2:      "shipping_city",
			wantQuery: []string{
				"1=1",
				"active = $1",
				"(first_name ILIKE '%' || $2::text || '%' OR shipping_city ILIKE '%' || $3::text || '%')",
			},
			wantArgs: []any{true, "João", "São Paulo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQuery, gotArgs := prepareOrLikeQuery(tt.filters1, tt.filters2, tt.initQuery, tt.initArgs, tt.key1, tt.key2)

			assert.Equal(t, tt.wantQuery, gotQuery)
			assert.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}
