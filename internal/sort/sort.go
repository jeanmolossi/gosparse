package sort

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Sorting é um custom type para os valores de ordenação aceitos
type Sorting byte

// Sort é um map de Sorting
//
// maps tem mais performance para busca por chaves
//
// Isso vai armazenar os valores aceitos para o parâmetro
// de busca "sort" no seguinte formato:
//
//	map[string]Sorting{
//		"created_at": DESC,
//		"title": ASC,
//	}
type Sort map[string]Sorting

// CtxKey é uma chave para o contexto.
// structs vazias tem mais performance.
type CtxKey struct{}

const (
	ASC Sorting = iota
	DESC

	SORT_PARAM string = "sort"
)

// extractSortFromQuery recebe a query e devolve um novo
// url.Values somente com a solicitação do parâmetro "sort"
func extractSortFromQuery(query url.Values) url.Values {
	extracted := url.Values{}

	if !query.Has(SORT_PARAM) {
		return extracted
	}

	extracted.Set(SORT_PARAM, query.Get(SORT_PARAM))
	return extracted
}

// Handle recebe um contexto e a query da request.
//
// Caso não haja o parâmetro de ordenação "sort" na query da request
// o próprio contexto será devolvido sem erro.
//
// Handle também aplica a validação de parâmetro a partir do Decode.
//
// Caso o parâmetro de "sort" não seja informado o servidor pode aplicar
// parâmetros de ordenação padrão.
func (s Sort) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	query = extractSortFromQuery(query)
	if len(query) == 0 {
		return ctx, nil
	}

	sort, err := Decode(query)
	if err != nil {
		return ctx, err
	}

	for field := range sort {
		if _, exists := s[field]; !exists {
			return ctx, fmt.Errorf("unsupported sorting by: %s", field)
		}
	}

	return context.WithValue(ctx, CtxKey{}, sort), nil
}

// hasInvalidChars checa se o campo recebido está dentro do range de
// caracteres aceitos de acordo com o regexp
func hasInvalidChars(f string) bool {
	return regexp.MustCompile(`[a-zA-Z_0-9]+`).ReplaceAllString(f, "") != ""
}

// Decode recebe a query e extrai os campos e valores da query.
func Decode(query url.Values) (Sort, error) {
	sort := Sort{}

	for _, field := range strings.Split(query.Get(SORT_PARAM), ",") {
		sorting := ASC

		if strings.HasPrefix(field, "-") {
			field = field[1:] // split "-" from starting string
			sorting = DESC
		}

		if hasInvalidChars(field) {
			return nil, fmt.Errorf("%s not acceptable, only [a-zA-Z_0-9]", field)
		}

		sort[field] = sorting
	}

	return sort, nil
}

// GetSort extraí o Sort recebido na request a partir do contexto.
func GetSort(ctx context.Context) (Sort, error) {
	if sort, present := ctx.Value(CtxKey{}).(Sort); present {
		return sort, nil
	}

	return nil, fmt.Errorf("sorter is not present on context")
}

// Get recebe o contexto e a chave do campo a qual deseja se recuperar
// a ordenação (ASC / DESC).
//
// Caso a chave do campo não esteja presente no contexto será devolvida
// a ordenação padrão ASC
func (s Sort) Get(ctx context.Context, f string) Sorting {
	if sorting, present := ctx.Value(CtxKey{}).(Sort); present {
		return sorting[f]
	}

	return ASC
}

// GetAll recebe o contexto e retorna o Sort contido no contexto.
//
// Caso o contexto não contenha o Sort, uma instância vazia será
// devolvida.
func (s Sort) GetAll(ctx context.Context) Sort {
	sort, err := GetSort(ctx)
	if err != nil {
		return make(Sort)
	}

	return sort
}

// AddField recebe o campo aceito no parâmetro "sort".
//
// Caso a chave recebida já esteja na lista de campos suportados, ela
// será ignorada.
//
// O sort por padrão de AddField é ASC, porém deve-se utilizar
// o valor armazenado em contexto
func (s Sort) AddField(field string) {
	if s == nil {
		s = make(Sort)
	}

	if _, duplicate := s[field]; !duplicate {
		s[field] = ASC
	}
}

// Options -----------------

type SortOpt func(*Sort)

// AcceptField é uma opção do construtor de *Sort. Essa função recebe um
// slice de strings que serão as chaves de campos aceitos ao validar a query da request
func AcceptField(fields ...string) SortOpt {
	return func(f *Sort) {
		if len(fields) == 0 {
			return
		}

		for _, field := range fields {
			f.AddField(field)
		}
	}
}

// Constructor -----------------

func New(opt ...SortOpt) *Sort {
	sort := &Sort{}

	for _, o := range opt {
		if o == nil {
			continue
		}

		o(sort)
	}

	return sort
}
