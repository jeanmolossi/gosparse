// Package sparsefieldsets
//
// Um cliente PODE solicitar que um endpoint retorne apenas campos específicos
// na resposta por tipo, incluindo um parâmetro "fields[TYPE]" de consulta.
//
// O valor de qualquer parâmetro "fields[TYPE]" DEVE ser uma lista separada
// por vírgula (U+002C COMMA, “,”) que se refere aos nomes dos campos a serem
// retornados. Um valor vazio indica que nenhum campo deve ser retornado.
//
// Se um cliente solicitar um conjunto restrito de campos para um determinado
// tipo de recurso, um endpoint NÃO DEVE incluir campos adicionais em objetos
// de recurso desse tipo em sua resposta.
//
// Se um cliente não especificar o conjunto de campos para um determinado tipo
// de recurso, o servidor PODE enviar todos os campos, um subconjunto de campos
// ou nenhum campo para esse tipo de recurso.
//
//	GET /articles?include=author&fields[articles]=title,body&fields[people]=name HTTP/1.1
//	Accept: application/hal+json
//
// Observação: o URI de exemplo acima mostra caracteres não codificados ([, ])
// simplesmente para facilitar a leitura. Na prática, esses caracteres devem ser
// codificados por porcentagem.
//
// Observação: esta seção se aplica a qualquer endpoint que responda com recursos
// como dados primários ou incluídos, independentemente do tipo de solicitação.
// Por exemplo, um servidor pode oferecer suporte a conjuntos de campos esparsos
// junto com uma solicitação "POST" para criar um recurso.
//
// # References
//
//   - https://docs.commercelayer.io/core/sparse-fieldsets
//   - https://jsonapi.org/format/#fetching-sparse-fieldsets
package sparsefieldsets
