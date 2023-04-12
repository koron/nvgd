# "trdsql" filter

[trdsql][orginal] filter filters CSV (or so) by SQL.

[trdsql query editor is provided][editor].
It will help to compose URL with trdsql filter.
See [tiny manual](#getting-started-with-editor) for editor tutorial.

[editor]:/trdsql/
[orginal]:https://github.com/noborus/trdsql

## Parameters

Name    | Type    | Requirements | Description
--------|---------|--------------|-------------
`q`     | String  | Mandatory    | SQL query string
`t`     | String  | Option       | override table name (default is `t`)
`ifmt`  | String  | Option       | input format (default is `CSV`)
`id`    | String  | Option       | input delimiter (default is `,`)
`ih`    | Boolean | Option       | input has header line or not (default is `false`)
`is`    | Integer | Option       | input skip lines (default is `0`) ???
`ir`    | Integer | Option       | input pre read lines (default is `1`) ???
`il`    | Boolean | Option       | input limit read lines (default is `false`)
`ijq`   | String  | Option       | input JQ query (default is empty) ???
`inull` | String  | Option       | input, NULL replacement string (default is none)
`ofmt`  | String  | Option       | output format (default is `CSV`)
`od`    | String  | Option       | output delimiter (default is `,`)
`oq`    | String  | Option       | output quote character (default is `"`)
`oaq`   | Boolean | Option       | output force quote all (default is `false`)
`ocrlf` | Boolean | Option       | output with CRLF (default is `false`)
`oh`    | Boolean | Option       | output append header line (default is `false`)
`onowrap` | Boolean | Option     | output without line wrap (default is `false`)
`onull` | String  | Option       | output, NULL replacement string (default is none)

### Input formats (`ifmt`)

Specify format of input.

* `CSV`
* `LTSV` - Labeled Tab-separataed Values
* `JSON`
* `TBLN` - https://tbln.dev/
* `TSV` - Tab-Separated Values format
* `PSV` - Pipe-Separated Values format

### Output formats (`ofmt`)

Specify format of output.

* `CSV`
* `LTSV` - Labeled Tab-separataed Values
* `JSON`
* `TBLN` - https://tbln.dev/
* `RAW`
* `MD` - Markdown
* `AT` - ASCII Table format
* `JSONL` - JSON Lines format(http://jsonlines.org/)

## Getting started with editor

The host and port address in the following text should be read according to
your environment.

1.  Open <http://127.0.0.1:9280/trdsql/>
2.  Input `http://127.0.0.1:9280/examples/line.csv` to **Source URL**
3.  Input `SELECT * FROM t` to a textarea at bottom of 2
4.  Push **Query** button.
    * This fills **Composed URL** and result area in **Output**.
    * You can share result of query with **Composed URL**
5.  You can change **Options** and push **Query** again to update **Output**
6.  Try `SELECT * FROM t WHERE date BETWEEN "2/" AND "3/"` as **Query**.
    Check changes of result.

## Notations

*   The query is processed by SQLite. You can use statements which supported by
    SQLite.  See <https://sqlite.org/lang_select.html> for details.
*   The table name is `t` as default, and you can change it by `t` parameter.
*   Column names are `c1`, `c2`, ... `c{n}` as default. But input header is
    used when `ih:true` parameter is specified.
*   `inull` and `onull` specify an alternate string for `NULL`.  The `inull`
    string is parsed as `NULL` in input.  `NULL` is replaced by the `onull`
    string in output.
*   SQLite works with memory (no disk). It will exhaust memory if try to treat
    very big input or heavy query.
*   The limitations in SQLite still apply
