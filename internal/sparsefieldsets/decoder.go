package sparsefieldsets

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// fieldMatcherSimple
	//
	// quando o campo NÃO tem o predicado a sequencia de match
	// segue a ordem como por exemplo:
	//
	//      fields[username]=john,anne
	//
	// Output
	//
	//  matches[0] = fields[username]
	//  matches[1] = username
	fieldMatcherSimple = regexp.MustCompile(`^fields\[([a-zA-Z_0-9]+)\]`).FindStringSubmatch
)

// extractField recebe a chave da querystring da
// request e extrai o nome do campo.
//
// para um formato inválido de chave, retorna uma string vazia e um erro
func extractField(f string) (string, error) {
	if f == SEARCH_PARAM {
		f = "fields[root]"
	}

	// se não houver match em um campo
	// o campo possui formato inválido
	matches := fieldMatcherSimple(f)

	if len(matches) == 0 {
		return "", fmt.Errorf("field has invalid format: %s", f)
	}

	return matches[1], nil

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
func Decode(query url.Values) (Fields, error) {
	fields := Fields{}

	for key, val := range query {
		field, err := extractField(key)
		if err != nil {
			return nil, err
		}

		fields[field] = resetValues(val)
	}

	return fields, nil
}
