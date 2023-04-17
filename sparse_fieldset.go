package gosparse

import (
	"strings"
)

type SparseFieldset struct {
	Include *Include
	Fields  *Fields
	Sort    *Sort
}

func (s *SparseFieldset) GetFields() []string {
	fields := make([]string, 0, len(s.Fields.FieldMap)+len(s.Sort.FieldMap)+len(s.Include.Fields()))

	isDuplicate := func(f string) bool {
		for _, field := range fields {
			if strings.EqualFold(field, f) {
				return true
			}
		}

		return false
	}

	for field := range s.Fields.FieldMap {
		if !isDuplicate(field) {
			fields = append(fields, field)
		}
	}

	for field := range s.Sort.FieldMap {
		if !isDuplicate(field) {
			fields = append(fields, field)
		}
	}

	for _, field := range s.Include.Fields() {
		if !isDuplicate(field) {
			fields = append(fields, field)
		}
	}

	return fields
}
