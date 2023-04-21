package pagination

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Pagination é um map de structs vazias
//
// maps tem mais performance para busca por chaves
//
// Isso vai armazenar os valores aceitos para o parâmetro
// de busca "page" no seguinte formato:
//
//	map[accepted]int{
//		"number": 1,
//		"size": 10,
//	}
type Pagination map[accepted]int

// PaginationOpt é uma assinatura para opções de configuração
// paga o construtor de Pagination
type PaginationOpt func(*Pagination)

// CtxKey é uma chage para o contexto.
// structs vazias tem mais performance.
type CtxKey struct{}

// accepted é um custom type para as propriedades aceitas
// para o parâmetro de busca "page"
type accepted string

const (
	PAGE_PARAM string = "page"

	SIZE   accepted = "size"
	NUMBER accepted = "number"
)

var (
	// pageMatcher
	//
	// A sequencia de match segue a ordem como por exemplo:
	//
	//      page[number]=10
	//
	// Output
	//
	//  matches[0] = page[number]
	//  matches[1] = number
	pageMatcher = regexp.MustCompile(`^page\[([a-zA-Z_0-9]+)\]$`).FindStringSubmatch
)

// extractPaginationFromQuery recebe a query e devolve um novo
// url.Values somente com a solicitação do parâmetro "page"
func extractPaginationFromQuery(query url.Values) url.Values {
	extracted := url.Values{}

	for k, v := range query {
		if strings.HasPrefix(k, PAGE_PARAM) {
			extracted[k] = v
		}
	}

	return extracted
}

// decode recebe a query e extrai os campos e valores da query.
func decode(query url.Values) (Pagination, error) {
	pagination := Pagination{}

	for key, val := range query {
		matches := pageMatcher(key)
		if len(matches) > 1 {
			k, err := StrToPageParam(matches[1])
			if err != nil {
				return pagination, err
			}

			// join para juntar quaisquer valores adicionais
			// exemplo: url.Values{"chave":{"1","2"}}
			value, err := strconv.Atoi(strings.Join(val, ""))
			if err != nil {
				return pagination, fmt.Errorf("pagination param %s should be int", k)
			}

			pagination[k] = value
			continue
		}

		return pagination, fmt.Errorf("missing prop on page param")
	}

	return pagination, nil
}

// StrToPageParam recebe a propriedade (size / number) e checa se
// é valida. Se for válida, retorna no formato accepted.
//
// Para propriedades inválidas retorna accepted vazio e um erro.
func StrToPageParam(p string) (accepted, error) {
	if strings.EqualFold(p, string(SIZE)) {
		return SIZE, nil
	}

	if strings.EqualFold(p, string(NUMBER)) {
		return NUMBER, nil
	}

	return "", fmt.Errorf("invalid pagination param %s", p)
}

// Handle recebe um contexto e a query da request.
//
// Caso não haja o parâmetro de paginação "page" na query da request
// o próprio contexto será devolvido sem erro.
//
// Handle também aplica a validação de parâmetro a partir do decode.
// Caso haja alguma propriedade de "page" que não é aceita, será
// retornado um erro de campo inválido / faltando.
//
// O parâmetro page só deve ser recebido com valores aceitos ou não deve ser utilizado.
func (p Pagination) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	query = extractPaginationFromQuery(query)
	if len(query) == 0 {
		return ctx, nil
	}

	pagination, err := decode(query)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, CtxKey{}, pagination), nil
}

// Get recebe o contexto e a chave do campo de "page" já validada e tratada.
//
// Caso o contexto não tenha o valor do campo e o campo não de match com
// nenhuma das propriedades aceitas (number / size), será retornado -1
//
// Exemplos:
//
//	// ...
//	pagination.Get(ctx, "invalid")  // -1
//	pagination.Get(ctx, SIZE)	// 10 - default
//	pagination.Get(ctx, NUMBER)	// 1  - default
func (p Pagination) Get(ctx context.Context, field accepted) int {
	if values, ok := ctx.Value(CtxKey{}).(Pagination); ok {
		return values[field]
	}

	switch field {
	case SIZE:
		return p[SIZE]
	case NUMBER:
		return p[NUMBER]
	}

	return -1
}

// DefaultPageSize altera o tamanho padrão do parâmetro size
//
// @Default = 10
//
// Caso o page size seja menor ou igual a zero, será definido o
// padrão @Default (10)
func DefaultPageSize(pageSize uint32) PaginationOpt {
	return func(p *Pagination) {
		if pageSize <= 0 {
			pageSize = 10
		}

		(*p)[SIZE] = int(pageSize)
	}
}

// Constructor -----------------

func New(opt ...PaginationOpt) *Pagination {
	pagination := &Pagination{
		NUMBER: 1,
		SIZE:   10,
	}

	if len(opt) == 0 {
		return pagination
	}

	for _, o := range opt {
		if o == nil {
			continue
		}

		o(pagination)
	}

	return pagination
}
