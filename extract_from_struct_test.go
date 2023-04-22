package gosparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Dummy struct {
	Title     string    `gosparse:"name:title;select;sort;filter"`
	CreatedAt time.Time `gosparse:"name:created_at;select;sort:desc;filter"`
	Nested    Nested    `gosparse:"name:nested;relation;"`
}

type Nested struct {
	Dummy bool `gosparse:"name:dummy;select;"`
}

func TestExtract(t *testing.T) {
	gosparse, err := Extract(Dummy{})

	require.Nil(t, err)
	require.NotEmpty(t, gosparse)
	require.NotNil(t, gosparse.Include)
	require.Len(t, gosparse.Include, 1) // only nested as relation
	require.NotNil(t, gosparse.Fieldset)
	require.Len(t, gosparse.Fieldset, 4) // contains all fields and root
	require.NotNil(t, gosparse.Filter)
	require.Len(t, gosparse.Filter, 2) // only fields with "filter" tag
	require.NotNil(t, gosparse.Pagination)
	require.Len(t, gosparse.Pagination, 2) // page number and page size
	require.NotNil(t, gosparse.Sort)
	require.Len(t, gosparse.Sort, 2) // only fields with "sort" tag
}

func TestExtractor(t *testing.T) {
	testtable := []struct {
		desc   string
		tag    string
		expect config
	}{
		{
			desc:   "should extract tag",
			tag:    `name:title`,
			expect: config{Name: "title"},
		},
		{
			desc: "should extract tag from anywhere",
			tag:  `select;sort;name:title;relation`,
			expect: config{
				Name:     "title",
				Select:   true,
				Sort:     true,
				Relation: true,
			},
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			have := extractor(tt.tag)
			require.EqualValues(t, tt.expect, have)
		})
	}
}
