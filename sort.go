package gosparse

import (
	"fmt"
	"regexp"
	"strings"
)

type SortDir int

type Sort struct {
	FieldMap map[string]SortDir
}

const (
	ASC SortDir = iota
	DESC
)

var (
	MinusMatch = regexp.
		MustCompile(`^-`).
		FindString
)

func (s *Sort) Add(values []string) error {
	if len(values) > 1 {
		return fmt.Errorf("sort values is invalid")
	}

	values = strings.Split(values[0], ",")

	for _, v := range values {
		if IsDesc(v) {
			s.FieldMap[v[1:]] = DESC
			continue
		}

		s.FieldMap[v] = ASC
	}

	return nil
}

func NewSorter() *Sort {
	return &Sort{
		FieldMap: make(map[string]SortDir),
	}
}

func IsDesc(value string) bool {
	return MinusMatch(value) != ""
}
