package sparsefieldsets

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
		expected Fields
		err      error
	}{
		{
			desc:  "should extract fields without predicate",
			query: url.Values{"fields[username]": {"john,anne"}},
			expected: Fields{
				"username": {"john", "anne"},
			},
		},
		{
			desc:  "should extract fields without predicate and join them",
			query: url.Values{"fields[username]": {"john,anne", "paul"}},
			expected: Fields{
				"username": {"john", "anne", "paul"},
			},
		},
		{
			desc:  "should extract fields with predicate",
			query: url.Values{"fields[username]": {"john", "anne"}},
			expected: Fields{
				"username": {"john", "anne"},
			},
		},
		{
			desc:     "should fail invalid format",
			query:    url.Values{"field-invalid": {}},
			expected: (Fields)(nil),
			err:      fmt.Errorf("field has invalid format: field-invalid"),
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			fields, err := Decode(tt.query)

			require.EqualValues(t, tt.err, err, "errors does not match")
			require.Equal(t, tt.expected, fields, "field map does not accepted")
		})
	}
}
