package filter

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// filterMatcherWithPredicate
	//
	// quando o campo tem o predicado a sequencia de match
	// segue a ordem como por exemplo:
	//
	//  filter[location_id_eq]=1
	//
	// Output:
	//
	//  matches[0] = filter[location_id_eq]
	//  matches[1] = location_id
	//  matches[2] = eq
	//
	// Se o campo não tiver predicado o regex não encontrará
	// nada e retornará nil
	filterMatcherWithPredicate = regexp.MustCompile(`^filter\[([a-zA-Z_0-9]+)_([a-zA-Z_0-9]+)\]`).FindStringSubmatch

	// filterMatcherSimple
	//
	// quando o campo NÃO tem o predicado a sequencia de match
	// segue a ordem como por exemplo:
	//
	//      filter[username]=john,anne
	//
	// Output
	//
	//  matches[0] = filter[username]
	//  matches[1] = username
	filterMatcherSimple = regexp.MustCompile(`^filter\[([a-zA-Z_0-9]+)\]`).FindStringSubmatch
)

// extractFilter recebe a chave da querystring da
// request e extrai o nome do campo e o predicado.
//
// para um formato inválido de chave, retorna uma string vazia, um NONE e um erro
func extractFilter(f string) (string, Predicate, error) {
	if f == SEARCH_PARAM {
		return "", NONE, fmt.Errorf("has no filter field param")
	}

	// se não houver match em um campo ele não possui
	// predicado ou então possui formato inválido
	matches := filterMatcherWithPredicate(f)

	if len(matches) == 0 {
		matches = filterMatcherSimple(f)

		// se o não houver match nem mesmo com campo simples
		// o campo possui formato inválido
		if len(matches) == 0 {
			return "", NONE, fmt.Errorf("filter has invalid format: %s", f)
		}

		return matches[1], NONE, nil
	}

	return matches[1], getPredicate(matches[2]), nil
}

// resetValues recebe o slice de strings da query e
// reseta os valores para um slice com os valores compilados
//
// Exemplos:
//
//	resetValues([]string{"john,anne", "paul"}) // []string{"john", "anne", "paul"}
//	resetValues([]string{"john,anne"}) // []string{"john", "anne"}
func resetValues(v []string) []string {
	joined := strings.Join(v, ",")
	v = strings.Split(joined, ",")

	return v
}

// Decode recebe a query e extrai os valores de campo e valores da query.
func Decode(query url.Values) (Filters, error) {
	fields := Filters{}

	for key, val := range query {
		field, predicate, err := extractFilter(key)
		if err != nil {
			return nil, err
		}

		fields[field] = Field{
			Predicate: predicate,
			Values:    resetValues(val),
		}
	}

	return fields, nil
}
