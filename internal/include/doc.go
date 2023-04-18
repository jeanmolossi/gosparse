// Package include
//
// Um endpoint PODE retornar recursos relacionados aos dados primários por padrão.
//
// Um endpoint também pode suportar um parâmetro "include" de consulta para
// permitir que o cliente personalize quais recursos relacionados devem ser
// retornados.
//
// Se um endpoint não suportar o parâmetro "include", ele DEVE responder
// 400 Bad Requesta qualquer solicitação que o inclua.
//
// Se um endpoint oferecer suporte ao parâmetro "include" e um cliente o fornecer:
//
//   - A resposta do servidor DEVE conter um documento composto,
//     mesmo que essa chave "included" contenha um array vazio
//     (porque os relacionamentos solicitados estão vazios).
//
//   - O servidor NÃO DEVE incluir objetos de recursos não solicitados
//     na seção "included" do documento composto.
//
// O valor do parâmetro "include" DEVE ser uma lista separada por
// vírgula (U+002C COMMA, “,”) de caminhos de relacionamento. Um caminho
// de relacionamento é uma lista separada por pontos (U+002E FULL-STOP, “.”) de
// nomes de relacionamento. Um valor vazio indica que nenhum recurso
// relacionado deve ser retornado.
//
// Se um servidor não conseguir identificar um caminho de relacionamento ou
// não suportar a inclusão de recursos de um caminho, ele DEVE responder com
// 400 Bad Request.
//
// Por exemplo, um caminho de relacionamento pode ser comments.author,
// onde "comments" é um relacionamento listado em um objeto "articles"
// e "author" é um relacionamento listado em um objeto "comments".
//
// Além disso, comentários podem ser solicitados com um artigo:
//
//	GET /articles/1?include=comments HTTP/1.1
//	Accept: application/hal+json
//
// Para solicitar recursos relacionados a outros recursos, um caminho separado
// por pontos para cada nome de relacionamento pode ser especificado:
//
//	GET /articles/1?include=comments.author HTTP/1.1
//	Accept: application/hal+json
//
// Observação: como os documentos compostos exigem vinculação completa (exceto
// quando a vinculação de relacionamento é excluída por conjuntos de campos
// esparsos), os recursos intermediários em um caminho de várias partes devem
// ser retornados junto com seus próprios nodes.
//
// Por exemplo, uma resposta a um pedido de "comments.author" deve incluir
// comments, bem como o "author" de cada um desses "comments".
//
// Observação: um servidor pode optar por expor um relacionamento profundamente
// aninhado, como "comments.author" um relacionamento direto com um nome alternativo,
// como "commentAuthors". Isso permitiria que um cliente solicitasse
// /articles/1?include=commentAuthors em vez de /articles/1?include=comments.author.
// Ao expor o relacionamento aninhado com um nome alternativo, o servidor ainda
// pode fornecer ligação completa em documentos compostos sem incluir recursos
// intermediários potencialmente indesejados.
//
// Vários recursos relacionados podem ser solicitados em uma lista separada por
// vírgulas:
//
//	GET /articles/1?include=comments.author,ratings HTTP/1.1
//	Accept: application/hal+json
//
// Além disso, recursos relacionados podem ser solicitados a partir de um endpoint
// de relacionamento:
//
//	GET /articles/1/relationships/comments?include=comments.author HTTP/1.1
//	Accept: application/hal+json
//
// Nesse caso, os dados primários seriam uma coleção de objetos identificadores
// de recursos que representam a vinculação aos comentários de um artigo,
// enquanto os comentários completos e os autores dos comentários seriam
// retornados como dados incluídos.
//
// Observação: esta seção se aplica a qualquer terminal que responda com dados
// primários, independentemente do tipo de solicitação.
// Por exemplo, um servidor pode suportar a inclusão de recursos relacionados
// junto com uma solicitação "POST" para criar um recurso ou relacionamento.
//
// # References
//
//   - https://docs.commercelayer.io/core/including-associations
//   - https://jsonapi.org/format/#fetching-includes
package include
