package sort

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	testtable := []struct {
		desc   string
		query  url.Values
		expect Sort
		err    error
	}{
		{
			desc:  "should recover sorting from fields",
			query: url.Values{"sort": {"-created_at,title"}},
			expect: Sort{
				"created_at": DESC,
				"title":      ASC,
			},
		},
		{
			desc:  "should fail if is unsupported field",
			query: url.Values{"sort": {"-created_at,title,posts"}},
			expect: Sort{
				"created_at": DESC,
				"title":      ASC,
				"posts":      ASC, // by default get works with this
			},
			err: fmt.Errorf("unsupported sorting by: posts"),
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			sort := New(AcceptField("created_at", "title"))

			ctx, err := sort.Handle(context.Background(), tt.query)
			require.EqualValues(t, err, tt.err)

			if tt.err != nil {
				return
			}

			for field, sorting := range tt.expect {
				require.Equal(t, sorting, sort.Get(ctx, field), "sorting does not match")
			}
		})
	}

	t.Run("should fail if value is unsupported", func(t *testing.T) {
		query := url.Values{"sort": {"!@#$%^"}}
		sort := New(AcceptField("!@#$%^"))

		ctx, err := sort.Handle(context.Background(), query)
		require.NotNil(t, err)
		require.EqualError(t, err, "!@#$%^ not acceptable, only [a-zA-Z_0-9]")

		require.EqualValues(t, ctx, context.Background())
	})
}

func TestGet(t *testing.T) {
	testtable := []struct {
		desc   string
		query  url.Values
		expect Sort
	}{
		{
			desc:   "should get sorting correctly",
			query:  url.Values{"sort": {"-created_at"}},
			expect: Sort{"created_at": DESC},
		},
		{
			desc:   "should get sorting default",
			query:  url.Values{},
			expect: Sort{"unknown": ASC},
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			sort := New(AcceptField("created_at", "title"))

			ctx, err := sort.Handle(context.Background(), tt.query)
			require.Nil(t, err)

			for field, sorting := range tt.expect {
				require.Equal(t, sorting, sort.Get(ctx, field))
			}
		})
	}
}

func TestGetSort(t *testing.T) {
	t.Run("should get context sort", func(t *testing.T) {
		query := url.Values{"sort": {"-created_at,title"}}
		sort := New(AcceptField("created_at", "title"))

		ctx, err := sort.Handle(context.Background(), query)
		require.Nil(t, err)

		recvSort, err := GetSort(ctx)
		require.Nil(t, err)
		require.EqualValues(t, recvSort["created_at"], DESC)
		require.EqualValues(t, recvSort["title"], ASC)
	})

	t.Run("should fail if has no sort on context", func(t *testing.T) {
		sort, err := GetSort(context.Background())
		require.Nil(t, sort)
		require.EqualError(t, err, "sorter is not present on context")
	})
}
