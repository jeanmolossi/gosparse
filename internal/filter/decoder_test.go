package filter

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	testtable := []struct {
		desc     string
		query    url.Values
		expected Filters
		err      error
	}{
		{
			desc:  "should extract filter without predicate",
			query: url.Values{"filter[username]": {"john,anne"}},
			expected: Filters{
				"username": Field{NONE, []string{"john", "anne"}},
			},
		},
		{
			desc:  "should extract filter without predicate and join them",
			query: url.Values{"filter[username]": {"john,anne", "paul"}},
			expected: Filters{
				"username": Field{NONE, []string{"john", "anne", "paul"}},
			},
		},
		{
			desc:     "should fail invalid format",
			query:    url.Values{"field-invalid": {}},
			expected: (Filters)(nil),
			err:      fmt.Errorf("filter has invalid format: field-invalid"),
		},

		{
			desc:  "should extract filter with predicate",
			query: url.Values{"filter[username_in]": {"john,anne"}},
			expected: Filters{
				"username": Field{IN, []string{"john", "anne"}},
			},
		},
		{
			desc:  "should extract filter with predicate and join them",
			query: url.Values{"filter[username_nin]": {"john,anne", "paul"}},
			expected: Filters{
				"username": Field{NIN, []string{"john", "anne", "paul"}},
			},
		},
		{
			desc:  "should extract filter with predicate",
			query: url.Values{"filter[username_eq]": {"john", "anne"}},
			expected: Filters{
				"username": Field{EQ, []string{"john", "anne"}},
			},
		},
		{
			desc:     "should fail filter without field param",
			query:    url.Values{"filter": {"anne"}},
			expected: (Filters)(nil),
			err:      fmt.Errorf("has no filter field param"),
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			filter, err := Decode(tt.query)

			require.EqualValues(t, tt.err, err, "errors does not match")
			require.Equal(t, tt.expected, filter, "field map does not accepted")
		})
	}
}
