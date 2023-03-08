# "trdsql" filter

[trdsql][trdsql] filter does ... filtering CSV/TSV (or so) by SQL.

[trdsql]:https://github.com/noborus/trdsql

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

### for output

Name   | Type    | Requirements | Description
-------|---------|--------------|-------------

### Input formats

* `CSV`
* `LTSV`
* `JSON`
* `TBLN`
* `TSV`
* `PSV`

### Output formats

* `CSV`
* `LTSV`
* `JSON`
* `TBLN`
* `RAW`
* `MD`
* `AT`
* `JSONL`
