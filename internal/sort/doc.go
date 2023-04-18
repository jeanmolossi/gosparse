// Package sort
//
// Um servidor pode optar por oferecer suporte a solicitações para classificar
// coleções de recursos de acordo com um ou mais critérios (“campos de classificação”).
//
// Observação: embora recomendado, os campos de classificação não precisam
// necessariamente corresponder aos nomes de atributo e relacionamento do recurso.
//
// Nota: Recomenda-se que os campos de classificação separados por pontos
// (U+002E FULL-STOP, “.”) sejam usados ​​para solicitar a classificação com base
// em atributos de relacionamento. Por exemplo, um campo de classificação
// "author.name" pode ser usado para solicitar que os dados primários sejam
// classificados com base no atributo de relacionamento "nome do author".
//
// Um endpoint PODE oferecer suporte a solicitações para classificar
// os dados primários com um parâmetro "sort" de consulta. O valor de "sort" DEVE
// representar campos de classificação.
//
//	GET /people?sort=age HTTP/1.1
//	Accept: application/hal+json
//
// Um endpoint PODE suportar vários campos de classificação, permitindo campos
// de classificação separados por vírgula (U+002C COMMA, “,”). Os campos de
// classificação DEVEM ser aplicados na ordem especificada.
//
//	GET /people?sort=age,name HTTP/1.1
//	Accept: application/hal+json
//
// A ordem de classificação para cada campo de classificação DEVE ser crescente,
// a menos que seja prefixada com um sinal de menos (U+002D HÍFEN-MENOS, “-“),
// caso em que DEVE ser decrescente.
//
//	GET /articles?sort=-created,title HTTP/1.1
//	Accept: application/hal+json
//
// O exemplo acima deve retornar os artigos mais recentes primeiro. Quaisquer
// artigos criados na mesma data serão classificados por seu título em ordem
// alfabética crescente.
//
// Se o servidor não suportar a classificação conforme especificado no parâmetro
// de consulta sort, ele DEVE retornar 400 Bad Request.
//
// Se a classificação for suportada pelo servidor e solicitada pelo cliente por
// meio do parâmetro de consulta "sort", o servidor DEVE retornar elementos da
// matriz de nível superior da resposta ordenada de acordo com os critérios
// especificados. O servidor PODE aplicar regras de classificação padrão ao
// nível superior se o parâmetro de solicitação "sort" não for especificado.
//
// Observação: esta seção se aplica a qualquer endpoint que responda com uma
// coleção de recursos como dados primários, independentemente do tipo de
// solicitação.
//
// # References
//
//   - https://docs.commercelayer.io/core/sorting-results
//   - https://jsonapi.org/format/#fetching-sorting
package sort
