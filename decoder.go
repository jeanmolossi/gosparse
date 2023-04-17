package gosparse

import (
	"net/url"
	"regexp"
)

var EmptyFieldset = SparseFieldset{
	Include: nil,
	Fields:  &Fields{},
	Sort:    NewSorter(),
}

func Bind(query url.Values) SparseFieldset {
	sparse := &SparseFieldset{
		Include: nil,
		Fields:  &Fields{},
		Sort:    NewSorter(),
	}

	for field, values := range query {
		if isInclude(field) {
			sparse.Include = NewInclude(values)
			continue
		}

		if isField(field) {
			sparse.Fields.Add(field, values)
			continue
		}

		if isSort(field) {
			sparse.Sort.Add(values)
			continue
		}
	}

	return *sparse
}

func isInclude(f string) bool {
	return regexp.
		MustCompile(`^include`).
		FindString(f) != ""
}

func isField(f string) bool {
	return regexp.
		MustCompile(`^fields`).
		FindString(f) != ""
}

func isSort(f string) bool {
	return regexp.
		MustCompile(`^sort`).
		FindString(f) != ""
}
