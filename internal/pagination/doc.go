// Package pagination
//
// Um servidor PODE escolher limitar o número de recursos retornados em uma
// resposta a um subconjunto (“página”) de todo o conjunto disponível.
//
// Um servidor PODE fornecer links para percorrer um conjunto de dados
// paginados (“links de paginação”).
//
// Os links de paginação DEVEM aparecer no objeto links que corresponde a uma
// coleção. Para paginar os dados primários, forneça links de paginação no
// objeto links de nível superior. Para paginar uma coleção incluída retornada
// em um documento composto, forneça links de paginação no objeto de links
// correspondente.
//
// As seguintes chaves DEVEM ser usadas para links de paginação:
//
//	first: a primeira página de dados
//	last: a última página de dados
//	prev: a página anterior de dados
//	next: a próxima página de dados
//
// As chaves DEVEM ser omitidas ou ter um valor "null" para indicar que um
// determinado link está indisponível.
//
// Os conceitos de ordem, conforme expressos na nomenclatura dos links de
// paginação, DEVEM permanecer consistentes com as regras de classificação do
// JSON:API.
//
// A família de parâmetros de consulta "page" é reservada para paginação.
// Servidores e clientes DEVEM usar esses parâmetros para operações de paginação.
//
// Observação: a API JSON é independente da estratégia de paginação usada por
// um servidor, mas a família de parâmetros "page" de consulta pode ser usada
// independentemente da estratégia empregada. Por exemplo, uma estratégia
// baseada em página pode usar parâmetros de consulta como page[number] e
// page[size], enquanto uma estratégia baseada em cursor pode usar page[cursor].
//
// Observação: esta seção se aplica a qualquer endpoint que responda com uma
// coleção de recursos como dados primários, independentemente do tipo de solicitação.
//
// # References
//
//   - https://docs.commercelayer.io/core/pagination
//   - https://jsonapi.org/format/#fetching-pagination
package pagination
