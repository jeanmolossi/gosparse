package sparsefieldsets

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldset(t *testing.T) {
	testtable := []struct {
		desc         string
		query        url.Values
		acceptable   []string
		expectations Fields
		err          error
	}{
		{
			desc: "should recover ok fields",
			query: url.Values{
				"other":            {},
				"fields":           {"root"},
				"fields[username]": {"john"},
				"fields[friends]":  {"anne,paul"},
			},
			acceptable: []string{"username", "friends"},
			expectations: Fields{
				"username": {"john"},
				"friends":  {"anne", "paul"},
			},
		},
		{
			desc:         "should return empty field",
			query:        url.Values{},
			acceptable:   []string{"username"},
			expectations: Fields{"unknown": {}},
		},
		{
			desc:         "should fail when received not accepted field",
			query:        url.Values{"fields[unknown]": {"any"}},
			expectations: Fields{},
			err:          fmt.Errorf("unsupported field resource: unknown"),
		},
	}

	for _, tt := range testtable {
		t.Run(tt.desc, func(t *testing.T) {
			fieldset := New(AcceptField(tt.acceptable...))

			ctx, err := fieldset.Handle(context.Background(), tt.query)
			require.EqualValues(t, err, tt.err)

			for key, values := range tt.expectations {
				field := fieldset.Get(ctx, key)
				require.Equal(t, values, field)
			}
		})
	}
}
