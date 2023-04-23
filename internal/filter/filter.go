package filter

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Filters é um map de structs vazias
//
// maps tem mais performance para busca por chaves
//
// Isso vai armazenar os valores aceitos para o parâmetro
// de busca "filter" no seguinte formato:
//
//	map[string]Field{
//		// ...
//		"comments": {},
//		"author": {},
//		"posts": {},
//		// ...
//	}
type Filters map[string]Field

// Field é a estrutura que armazena a configuração
// de um campo específico no parâmetro de busca "filter"
type Field struct {
	Predicate Predicate
	Values    []string
}

// FiltersOpt é umas assinatura para opções de configuração
// para o construtor de Filters
type FiltersOpt func(*Filters)

// CtxKey é uma chave para o contexto.
// struct vazias são mais performaticas até mesmo que strings
type CtxKey struct{}

const (
	SEARCH_PARAM string = "filter"
)

// extractFilterFromQuery recebe a query e devolve um novo
// url.Values somente com a solicitação do parâmetro "filter"
func extractFilterFromQuery(query url.Values) url.Values {
	extracted := url.Values{}

	for k, v := range query {
		if strings.HasPrefix(k, SEARCH_PARAM) {
			extracted[k] = v
		}
	}

	return extracted
}

// Handle recebe um contexto e a query da request.
//
// Caso não haja o parâmetro de busca "filter" na query da request o
// próprio contexto será devolvido sem erro.
//
// Handle também aplica a validação de parâmetro "filter" que é aceito.
// Caso haja algum valor de "filter" que não está definido como "AcceptFilter"
// será retornado um erro de recurso de campo não suportado, uma vez que o
// parâmetro "filter" só deve ser recebido com valores aceitos ou não deve ser utilizado.
func (f Filters) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	query = extractFilterFromQuery(query)
	if len(query) == 0 {
		return ctx, nil
	}

	filters, err := Decode(query)
	if err != nil {
		return ctx, err
	}

	for filter := range filters {
		if _, exists := f[filter]; !exists {
			return ctx, fmt.Errorf("unsupported filter resource: %s", filter)
		}
	}

	return context.WithValue(ctx, CtxKey{}, filters), nil
}

// Get recebe o contexto e a chave do campo de "filter" já validado e tratado.
//
// Caso o contexto não tenha o valor do campo será retornado um Field zero valued.
func (f Filters) Get(ctx context.Context, field string) Field {
	if values, ok := ctx.Value(CtxKey{}).(Filters); ok {
		return values[field]
	}

	return Field{}
}

// GetAll recebe o contexto e retorna uma instância completa
// dos filtros contidos no contexto.
//
// Caso os filtros não estejam presentes no contexto, será
// devolvido uma instância vazia dos Filters
func (f Filters) GetAll(ctx context.Context) Filters {
	if values, ok := ctx.Value(CtxKey{}).(Filters); ok {
		return values
	}

	return make(Filters)
}

// AddFilter recebe a chave do campo aceito no parâmetro "filter".
//
// Caso a chave recebida já esteja na lista de campos suportados, ela
// será ignorada.
func (f Filters) AddFilter(filter string) {
	if f == nil {
		f = make(Filters)
	}

	if _, duplicate := f[filter]; !duplicate {
		f[filter] = Field{}
	}
}

// Options -----------------

// AcceptField é uma opção do construtor de *FiltersOpt. Essa função recebe um
// slice de strings que serão as chaves de campos aceitos ao validar a query da request
func AcceptField(filters ...string) FiltersOpt {
	return func(f *Filters) {
		if len(filters) == 0 {
			return
		}

		for _, field := range filters {
			f.AddFilter(field)
		}
	}
}

// Constructor -----------------

func New(opt ...FiltersOpt) *Filters {
	fields := &Filters{}

	for _, o := range opt {
		if o == nil {
			continue
		}

		o(fields)
	}

	return fields
}
