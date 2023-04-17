package gosparse

import (
	"regexp"
	"strings"
)

type Action int

type Field struct {
	Predicate Predicate
	Action    Action
}

type Fields struct {
	FieldMap map[string]Field
}

const (
	LIST Action = iota
	HIDE
)

var (
	// Regex:
	//
	// 	Should start with "fields["
	//	"^field\["
	//
	//	Should has letters and numbers between a underscore
	// 	"([a-zA-Z_0-9]+)_([a-zA-Z_0-9]+)"
	//
	// 	Should end with "]"
	//	"\]$"
	//
	// It will match fiels with following format:
	//
	// 	fields[fieldname_predicate]
	FieldMatcher = regexp.
		MustCompile(`^fields\[([a-zA-Z_0-9]+)_([a-zA-Z_0-9]+)\]$`).
		FindStringSubmatch
)

func resetValues(v []string) []string {
	if len(v) == 0 {
		return []string{""}
	}

	if len(v) == 1 {
		v = strings.Split(v[0], ",")
	}

	return v
}

func getAction(v string) Action {
	action := regexp.
		MustCompile(`^-`).
		FindString

	if action(v) == "" {
		return HIDE
	}

	return LIST
}

func removeMinus(v string) string {
	return regexp.
		MustCompile(`^-?`).
		ReplaceAllString(v, "")
}

func extractFieldValues(v []string) Action {
	v = resetValues(v)

	for _, value := range v {
		return getAction(value)
	}

	return LIST
}

func GetFieldName(field string) string {
	matches := FieldMatcher(field)
	if len(matches) >= 1 {
		return matches[1]
	}

	isRoot := regexp.
		MustCompile(`^fields(\[document\])?$`).
		FindString(field) != ""

	if isRoot {
		return DOC
	}

	return ""
}

func (f *Fields) Add(field string, values []string) {
	if f.FieldMap == nil {
		f.FieldMap = make(map[string]Field)
	}

	if field == "" || field == "fields" {
		field = "fields[document]"
	}

	fieldname := GetFieldName(field)

	if fieldname == "document" {
		values = resetValues(values)
		for _, value := range values {
			field := removeMinus(value)
			if _, duplicate := f.FieldMap[field]; duplicate {
				continue
			}

			f.FieldMap[field] = Field{
				Predicate: NONE,
				Action:    getAction(value),
			}
		}

		return
	}

	f.FieldMap[GetFieldName(field)] = Field{
		Predicate: GetPredicate(field),
		Action:    extractFieldValues(values),
	}
}

func (f *Fields) Get(field string) Field {
	return f.FieldMap[field]
}

func (f *Field) Get() Action {
	return f.Action
}