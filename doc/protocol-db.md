# DB protocol

## Multiple queries

* Multiple queries are accepted as a request.
* Queries are separated by `;` at end of line.
* The reponse includes only the last query's one. It sould be `SELECT` or so.
* Queries before the last should be used to modify session parameters by `SET` or so.
* No limitation on count of queries in a request, but those are executed in a transaction, modifications are temporary.
* There are only tiny and simple syntax check, there are no warranty for corner cases.

日本語訳 (Japanese translation):

* クエリに複数文を指定できるようにしました。
* 1つのクエリは行末の `;` で識別・分割されます。
* 最後のクエリはその出力がレスポンスとして返されるので、 `SELECT` である必要があります。
* 最後以外のクエリは `SET` でセッションパラメーターを変更するような用途を想定しています。
* 1回に実行できるクエリ数に上限はありませんが、1つのトランザクション内で実行されるため変更したセッションパラメーターは、その実行限りのものです。
* 構文チェックは簡易のモノであるため、コーナーケースについては保証できません。

## Notations

* `max_rows` is applied for `SELECT` queries without `COUNT` nor `LIMIT` clauses.
* MySQL's connection is not pooled, to reset variables which changed in a "query" request.
