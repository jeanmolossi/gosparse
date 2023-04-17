package gosparse_test

import (
	"gosparse"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type DummyPosts struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Published bool   `json:"published"`
}

type Dummy struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Posts     []DummyPosts `json:"posts"`
	Followers []Dummy      `json:"followers"`
	CreatedAt time.Time    `json:"createdAt"`
}

func TestBind(t *testing.T) {
	qry := url.Values{
		"include":             {"document,posts,followers"},
		"fields[username_in]": {"john,anne"},
		"fields[document]":    {"-name,username"},
		"sort":                {"-createdAt"},
	}

	wantFields := []string{"name", "username", "createdAt", "posts", "followers"}
	sparse := gosparse.Bind(qry)

	require.NotEqualValues(t, sparse, gosparse.EmptyFieldset)
	require.ElementsMatch(t, wantFields, sparse.GetFields())
	require.Equal(t, sparse.Fields.Get("username").Predicate, gosparse.IN)

	t.Logf(
		"\nTARGET RESULT:\nINCLUDE: %+v\nFIELDS: %+v\nSORT: %+v\n",
		sparse.Include,
		sparse.Fields,
		sparse.Sort,
	)
}
