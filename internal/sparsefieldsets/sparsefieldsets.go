package sparsefieldsets

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Fieldset é um map de structs vazias
//
// maps tem mais performance para busca por chaves
//
// Isso vai armazenar os valores aceitos para o parâmetro
// de busca "fields" no seguinte formato:
//
//	map[string]struct{}{
//		// ...
//		"comments": {},
//		"author": {},
//		"posts": {},
//		// ...
//	}
type Fieldset map[string]struct{}

// Fields é o map que irá salvar a solicitação de
// campo com predicado e valores
type Fields map[string][]string

// FieldsetOpt é uma assinatura para opções de configuração
// para o construtor de Fieldset
type FieldsetOpt func(*Fieldset)

type CtxKey struct{}

const (
	SEARCH_PARAM string = "fields"
)

// extractFieldFromQuery recebe a query e devolve um novo
// url.Values somente com a solicitação do parâmetro "fields"
func extractFieldFromQuery(query url.Values) url.Values {
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
// Caso não haja o parâmetro de busca "fields" na query da request o
// próprio contexto será devolvido sem erro.
//
// Handle também aplica a validação de parâmetro "fields" que é aceito.
// Caso haja algum valor de "fields" que não está definido como "AcceptField"
// será retornado um erro de recurso de campo não suportado, uma vez que o
// parâmetro "fields" só deve ser recebido com valores aceitos ou não deve ser utilizado.
func (f Fieldset) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	query = extractFieldFromQuery(query)
	if len(query) == 0 {
		return ctx, nil
	}

	fields, err := Decode(query)
	if err != nil {
		return ctx, err
	}

	for field := range fields {
		if _, exists := f[field]; !exists {
			return ctx, fmt.Errorf("unsupported field resource: %s", field)
		}
	}

	return context.WithValue(ctx, CtxKey{}, fields), nil
}

// Get recebe o contexto e a chave do campo de "fields" já validado e tratado.
//
// Caso o contexto não tenha o valor do campo será retornado um slice vazio.
func (f Fieldset) Get(ctx context.Context, field string) []string {
	if values, ok := ctx.Value(CtxKey{}).(Fields); ok {
		return values[field]
	}

	return make([]string, 0)
}

// GetAll recebe o contexto e devolve os fields contidos.
//
// Caso não haja Fields no contexto será devolvida uma instância vazia
func (f Fieldset) GetAll(ctx context.Context) Fields {
	if values, ok := ctx.Value(CtxKey{}).(Fields); ok {
		return values
	}

	return make(Fields)
}

// AddField recebe a chave do campo aceito no parâmetro "fields".
//
// Caso a chave recebida já esteja na lista de campos suportados, ela
// será ignorada.
func (f Fieldset) AddField(field string) {
	if f == nil {
		f = make(Fieldset)
	}

	if _, duplicate := f[field]; !duplicate {
		f[field] = struct{}{}
	}
}

// Options -----------------

// AcceptField é uma opção do construtor de *Fieldset. Essa função recebe um
// slice de strings que serão as chaves de campos aceitos ao validar a query da request
func AcceptField(fields ...string) FieldsetOpt {
	return func(f *Fieldset) {
		if len(fields) == 0 {
			return
		}

		for _, field := range fields {
			f.AddField(field)
		}
	}
}

// Constructor -----------------

func New(opt ...FieldsetOpt) *Fieldset {
	fields := &Fieldset{}

	for _, o := range opt {
		if o == nil {
			continue
		}

		o(fields)
	}

	AcceptField("root")(fields)

	return fields
}
