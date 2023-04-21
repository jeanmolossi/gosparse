package pagination

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPagination(t *testing.T) {
	testtable := []struct {
		desc         string
		query        url.Values
		expectNumber int
		expectSize   int
		err          error
	}{
		{
			desc: "should get correct params",
			query: url.Values{
				"page[number]": {"1"},
				"page[size]":   {"30"},
			},
			expectNumber: 1,
			expectSize:   30,
		},
		{
			desc:         "should fail if passed invalid prop",
			query:        url.Values{"page[total]": {"1"}},
			expectNumber: 1,
			expectSize:   10,
			err:          fmt.Errorf("invalid pagination param total"),
		},
		{
			desc:         "should fail if can not parse to int",
			query:        url.Values{"page[number]": {"a"}},
			expectNumber: 1,
			expectSize:   10,
			err:          fmt.Errorf("pagination param number should be int"),
		},
		{
			desc:         "should ignore handle if has no page param",
			query:        url.Values{},
			expectNumber: 1,
			expectSize:   10,
			err:          nil,
		},
		{
			desc:         "should fail if not specify page prop param",
			query:        url.Values{"page": {"1"}}, // expected to be page[number]
			expectNumber: 1,
			expectSize:   10,
			err:          fmt.Errorf("missing prop on page param"),
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			pagination := New()
			ctx, err := pagination.Handle(
				context.Background(),
				tt.query,
			)

			require.EqualValues(t, err, tt.err)
			require.Equal(t, tt.expectSize, pagination.Get(ctx, SIZE))
			require.Equal(t, tt.expectNumber, pagination.Get(ctx, NUMBER))
		})
	}

	t.Run("should instantiate with another default", func(t *testing.T) {
		pagination := New(DefaultPageSize(30))
		ctx, err := pagination.Handle(
			context.Background(),
			url.Values{},
		)

		require.Nil(t, err)
		require.Equal(t, 30, pagination.Get(ctx, SIZE))
	})

	t.Run("should instantiate with another default", func(t *testing.T) {
		pagination := New(DefaultPageSize(0))
		ctx, err := pagination.Handle(
			context.Background(),
			url.Values{},
		)

		require.Nil(t, err)
		require.Equal(t, 10, pagination.Get(ctx, SIZE))
	})

	t.Run("should instantiate with default", func(t *testing.T) {
		pagination := New(nil)
		ctx, err := pagination.Handle(
			context.Background(),
			url.Values{},
		)

		require.Nil(t, err)
		require.Equal(t, 10, pagination.Get(ctx, SIZE))
	})
}
