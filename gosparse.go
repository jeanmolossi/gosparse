package gosparse

import (
	"context"
	"net/url"

	"github.com/jeanmolossi/gosparse/internal/filter"
	"github.com/jeanmolossi/gosparse/internal/include"
	"github.com/jeanmolossi/gosparse/internal/pagination"
	"github.com/jeanmolossi/gosparse/internal/sort"
	"github.com/jeanmolossi/gosparse/internal/sparsefieldsets"
)

// Gosparse contém a configuração base de parâmetros aceitos
type Gosparse struct {
	Include    include.Includes
	Fieldset   sparsefieldsets.Fieldset
	Filter     filter.Filters
	Pagination pagination.Pagination
	Sort       sort.Sort
}

// Handle recebe o contexto e a querystring da request e
// extraí todos os parâmetros de filtro, seleção e ordenação.
//
//   - Include
//   - Fieldset
//   - Filter
//   - Pagination
//   - Sort
func (g Gosparse) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	ctx, err := g.Include.Handle(ctx, query)
	if err != nil {
		return ctx, err
	}

	ctx, err = g.Fieldset.Handle(ctx, query)
	if err != nil {
		return ctx, err
	}

	ctx, err = g.Filter.Handle(ctx, query)
	if err != nil {
		return ctx, err
	}

	ctx, err = g.Pagination.Handle(ctx, query)
	if err != nil {
		return ctx, err
	}

	ctx, err = g.Sort.Handle(ctx, query)
	if err != nil {
		return ctx, err
	}

	return ctx, err
}

// Options ------------------------------

type GosparseOpt func(*Gosparse)

func AcceptRelations(rels ...string) GosparseOpt {
	return func(g *Gosparse) {
		if g.Include == nil {
			g.Include = *include.New(include.AcceptRel(rels...))
			return
		}

		for _, rel := range rels {
			g.Include.AddRel(rel)
		}
	}
}

func AcceptFields(fields ...string) GosparseOpt {
	return func(g *Gosparse) {
		if g.Fieldset == nil {
			g.Fieldset = *sparsefieldsets.New(sparsefieldsets.AcceptField(fields...))
			return
		}

		for _, field := range fields {
			g.Fieldset.AddField(field)
		}
	}
}

func AcceptFilters(filters ...string) GosparseOpt {
	return func(g *Gosparse) {
		if g.Filter == nil {
			g.Filter = *filter.New(filter.AcceptField(filters...))
			return
		}

		for _, filter := range filters {
			g.Filter.AddFilter(filter)
		}
	}
}

func AcceptPagination(size uint32) GosparseOpt {
	return func(g *Gosparse) {
		if g.Pagination == nil {
			g.Pagination = *pagination.New(pagination.DefaultPageSize(size))
			return
		}

		g.Pagination[pagination.SIZE] = int(size)
	}
}

func AcceptSortBy(fields ...string) GosparseOpt {
	return func(g *Gosparse) {
		if g.Sort == nil {
			g.Sort = *sort.New(sort.AcceptField(fields...))
			return
		}

		for _, field := range fields {
			g.Sort.AddField(field)
		}
	}
}

// Constructor --------------------------

func New(options ...GosparseOpt) Gosparse {
	gosparse := Gosparse{}

	for _, opt := range options {
		if opt == nil {
			continue
		}

		opt(&gosparse)
	}

	return gosparse
}
