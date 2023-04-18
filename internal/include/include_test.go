package include_test

import (
	"context"
	"fmt"
	"gosparse/internal/include"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	require.NotNil(t, include.New(nil))
}

func TestHandle(t *testing.T) {
	testtable := []struct {
		testdescription string
		include         string
		acceptable      []string
		want            []string
		err             error
	}{
		{
			testdescription: "should recover include values",
			include:         "comments,comments.author",
			acceptable:      []string{"comments", "comments.author"},
			want:            []string{"comments", "comments.author"},
			err:             nil,
		},
		{
			testdescription: "should recover include values",
			include:         "",
			acceptable:      []string{"comments", "comments.author"},
			want:            []string{},
			err:             nil,
		},
		{
			testdescription: "should fail if receive non acceptable value",
			include:         "posts",
			acceptable:      []string{"comments"},
			want:            []string{},
			err:             fmt.Errorf("unsupported include relation posts"),
		},
		{
			testdescription: "should try with wrong values",
			include:         "",
			acceptable:      nil,
			want:            []string{},
			err:             nil,
		},
	}

	for _, tt := range testtable {
		t.Run(tt.testdescription, func(t *testing.T) {
			query := url.Values{include.SEARCH_PARAM: {tt.include}}

			inc := include.New(include.AcceptRel(tt.acceptable...))

			ctx, err := inc.Handle(context.Background(), query)
			require.EqualValues(t, err, tt.err)

			rel := inc.Get(ctx)
			require.ElementsMatch(t, tt.want, rel)
		})
	}

}
