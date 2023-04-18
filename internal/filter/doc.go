// Package filter
//
// Uma “família de parâmetros de consulta” é o conjunto de todos os parâmetros
// de consulta cujo nome começa com um “nome de base”, seguido por zero ou mais
// instâncias de colchetes vazios, ou seja, nomes de membros entre
// colchetes. A família é referida pelo seu nome base.
//
// Por exemplo, a família de parâmetros de consulta inclui parâmetros "filter"
// denominados: filter, filter[x], filter[], etc. No entanto, não é um nome
// de parâmetro válido na família, porque não é um nome de membro válido.
//
//	`filter[x][] filter[][] filter[x][y] filter[_]_`
//
// # Parâmetros de consulta específicos da extensão
//
// O nome base de cada parâmetro de consulta introduzido por uma extensão DEVE
// ser prefixado com o namespace da extensão seguido por dois pontos (:).
// O restante do nome base DEVE conter apenas os caracteres [az] (U+0061 a U+007A, “az”).
//
// Parâmetros de consulta específicos da implementação
// As implementações PODEM oferecer suporte a parâmetros de consulta personalizados.
// No entanto, os nomes desses parâmetros de consulta DEVEM vir de uma família
// cujo nome base seja um nome de membro legal e também contenha pelo menos um
// caractere não az (ou seja, fora de U+0061 a U+007A).
//
// É RECOMENDADO que uma letra maiúscula (por exemplo, camelCasing) seja usada
// para atender ao requisito acima.
//
// Se um servidor encontrar um parâmetro de consulta que não segue as convenções
// de nomenclatura acima, ou o servidor não souber como processá-lo como um
// parâmetro de consulta desta especificação, ele DEVE retornar 400 Bad Request.
//
// Observação: ao proibir o uso de parâmetros de consulta que contenham apenas
// os caracteres [az], o JSON:API reserva a capacidade de padronizar parâmetros
// de consulta adicionais posteriormente sem entrar em conflito com as implementações existentes.
//
// # References
//
//   - https://docs.commercelayer.io/core/filtering-data
//   - https://jsonapi.org/format/#fetching-filtering
//   - https://jsonapi.org/format/#query-parameters-families
package filter
