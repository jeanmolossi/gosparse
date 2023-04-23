package gosparse

import (
	"context"
	"net/url"
	"testing"

	"github.com/jeanmolossi/gosparse/internal/filter"
	"github.com/jeanmolossi/gosparse/internal/pagination"
	"github.com/jeanmolossi/gosparse/internal/sort"
	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	query := url.Values{
		"include":                  {"nested"},
		"fields":                   {"title"},
		"fields[nested]":           {"dummy"},
		"filter[created_at_start]": {"2023-01-01"},
		"page[number]":             {"1"},
		"page[size]":               {"15"},
		"sort":                     {"-created_at"},
	}

	gosparse, err := Extract(Dummy{})
	require.Nil(t, err)

	ctx, err := gosparse.Handle(context.Background(), query)
	require.Nil(t, err)
	require.NotNil(t, ctx)

	// Include assertions
	require.Contains(t, gosparse.Include.Get(ctx), "nested")

	// Fieldset assertions
	require.Contains(t, gosparse.Fieldset.GetAll(ctx), "nested")
	require.Contains(t, gosparse.Fieldset.GetAll(ctx), "root")

	// Filter assertions
	require.Contains(t, gosparse.Filter.GetAll(ctx), "created_at")
	require.EqualValues(t,
		filter.Field{Predicate: filter.START, Values: []string{"2023-01-01"}},
		gosparse.Filter.Get(ctx, "created_at"),
	)

	// Pagination assertions
	require.Equal(t, 1, gosparse.Pagination.Get(ctx, pagination.NUMBER))
	require.Equal(t, 15, gosparse.Pagination.Get(ctx, pagination.SIZE))

	// Sort assertions
	require.Equal(t, sort.DESC, gosparse.Sort.Get(ctx, "created_at"))
}
