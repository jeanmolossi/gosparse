package gosparse

import (
	"strings"
)

type Include struct {
	root bool
	doc  []string
}

const (
	DOC = "document"
)

func (i *Include) Fields() []string {
	total := len(i.doc)
	if i.root {
		total = total - 1
	}

	f := make([]string, 0, total)

	for _, field := range i.doc {
		if field == DOC {
			continue
		}

		f = append(f, field)
	}

	return f
}

func NewInclude(inc []string) *Include {
	root := false
	if len(inc) == 0 {
		inc = append(inc, DOC)
		root = true
	}

	if len(inc) == 1 {
		inc = strings.Split(inc[0], ",")
	}

	notEmptyFields := make([]string, 0, len(inc))

	if !root {
		for _, field := range inc {
			if field == "" {
				continue
			}

			if field == DOC {
				root = true
			}

			notEmptyFields = append(notEmptyFields, field)
		}
	}

	if len(notEmptyFields) == 0 {
		notEmptyFields = append(notEmptyFields, DOC)
		root = true
	}

	return &Include{
		root: root,
		doc:  notEmptyFields,
	}
}
