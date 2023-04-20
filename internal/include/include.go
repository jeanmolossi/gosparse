package include

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Includes é um map de structs vazias
//
// maps tem mais performance para busca por chaves
//
// Isso vai armazenar os valores aceitos para o parâmetro
// de busca "include" no seguinte formato:
//
//	map[string]struct{}{
//		// ...
//		"comments": {},
//		"comments.author": {},
//		"posts": {},
//		// ...
//	}
type Includes map[string]struct{}

// IncludeOpt é uma assinatura para opções de configuração
// para o construtor de Includes
type IncludeOpt func(*Includes)

type CtxKey struct{}

const (
	SEARCH_PARAM string = "include"
)

// Handle vai receber um contexto e a query da request.
//
// Caso não haja o parâmetro de busca "include" na query da request o
// próprio contexto será devolvido
//
// Handle também aplica a validação de parâmetro "include" que é aceito.
// Caso haja algum valor de include que não está definido como "AcceptRel"
// será retornado um erro de relação não suportada, uma vez que o parâmetro
// "include" só deve ser recebido com valores aceitos ou não deve ser utilizado.
//
// É importante ressaltar que o contexto será sobrescrito por um novo contexto,
// agora com o valor do "include", portanto, repasse o contexto retornado de
// Handle para as chamadas seguintes que poderão recuperar os valores de "include".
func (r Includes) Handle(ctx context.Context, query url.Values) (context.Context, error) {
	if !query.Has(SEARCH_PARAM) || query.Get(SEARCH_PARAM) == "" {
		return ctx, nil
	}

	values := strings.Split(query[SEARCH_PARAM][0], ",")

	for _, val := range values {
		if _, exists := r[val]; !exists {
			return ctx, fmt.Errorf("unsupported include relation %s", val)
		}
	}

	return context.WithValue(ctx, CtxKey{}, values), nil
}

// Get recebe o contexto com os valores de "include" já validados e tratados.
//
// Caso o contexto não tenha os valores de "include" será retornado um slice
// de strings vazio.
//
//	[]string{}
func (r Includes) Get(ctx context.Context) []string {
	if values, ok := ctx.Value(CtxKey{}).([]string); ok {
		return values
	}

	return make([]string, 0)
}

// AddRel recebe uma nova relação de campos aceitos no parâmetro "include".
//
// Caso a relação recebida seja duplicada, ela não será adicionada às relações
// aceitas no campo
func (r Includes) AddRel(rel string) {
	if r == nil {
		r = make(Includes)
	}

	if _, duplicate := r[rel]; !duplicate {
		r[rel] = struct{}{}
	}
}

// Options -----------------

// AcceptRel é uma opção do construtor de *Includes. Essa função recebe um
// slice de strings que serão as relações aceitas ao validar a query da request
func AcceptRel(rels ...string) IncludeOpt {
	return func(i *Includes) {
		if len(rels) == 0 {
			return
		}

		for _, rel := range rels {
			i.AddRel(rel)
		}
	}
}

// Constructor -----------------

func New(opt ...IncludeOpt) *Includes {
	inc := &Includes{}

	for _, o := range opt {
		if o == nil {
			continue
		}

		o(inc)
	}

	return inc
}
