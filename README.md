# NVGD - Night Vision Goggles Daemon

HTTP file server to help DevOps.

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/nvgd)](https://pkg.go.dev/github.com/koron/nvgd)
[![Actions/Go](https://github.com/koron/nvgd/workflows/Go/badge.svg)](https://github.com/koron/nvgd/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/nvgd)](https://goreportcard.com/report/github.com/koron/nvgd)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/koron/nvgd)

NVGD (Night Vision Goggles Daemon) is an HTTP file server designed to help
DevOps professionals access and transform data from various sources.  It
provides a unified interface to retrieve data from local files, databases, AWS
S3, Redis, and other sources, then apply transformations like filtering,
formatting, and visualization before serving the results to clients.

Index:

* [How to use](#how-to-use)
* [Acceptable path](#acceptable-path)
    * [Protocols](#protocols)
* [Configuration file](#configuration-file)
    * [File Protocol Handler](#file-protocol-handler)
    * [Command Protocol Handlers](#command-protocol-handlers)
    * [S3 Protocol Handlers](#config-s3-protocol-handlers)
    * [Config DB Protocol Handler](#config-db-protocol-handler)
    * [Configure filters](#configure-filters)
    * [Default Filters](#default-filters)
* [Filters](#filters)
* [Prefix Aliases](#prefix-aliases)

## How to use

Download an archive file which matches to your environment from [latest
release](https://github.com/koron/nvgd/releases/latest).

Extract an executable `nvgd` or `nvgd.exe` from the archive file.  Then copy it
to one of directory in PATH environment variable. (ex. `/usr/local/bin`)

Run:

    $ nvgd

Access:

    $ curl http://127.0.0.1:9280/file:///var/log/message/httpd.log?tail=limit:25

NOTE: Pre-compiled binary for Linux is built with newer glibc. So it can't be
run on Linux with old glibc, like CentOS 7 or so.  In that case, you must
compile nvgd by your self. Please check next section to build from source.

### Build from source

Requirements to build:

* Go 1.24.0 or above (1.24.3 is recommended)
    * CGO enabled

How to install or upgrade

```console
$ go install github.com/koron/nvgd@latest
```

See also: [How to build on CentOS 7](doc/build-centos7.md)

## Acceptable path

Nvgd accepts path in these like format:

    /{protocol}://{args/for/protocol}[?{filters}]

### Protocols

Nvgd supports these `protocol`s:

* `file` - `/file:///path/to/source`
    * support [glob][globspec] like `*`

        ```
        /files:///var/log/access*.log
        ```

* `command` - result of pre-defined commands
* `s3obj`
    * get object: `/s3obj://bucket-name/key/to/object`
* `s3list`
    * list common prefixes and objects: `/s3list://bucket-name/prefix/of/key`
* `db` - query pre-defined databases
    * query `id` and `email` form users in `db_pq`:

        ```
        /db://db_pq/select id,email from users
        ```

    * support multiple databases:

        ```
        /db://db_pq2/foo/select id,email from users
        /db://db_pq2/bar/select id,email from users
        ```

        This searches from `foo` and `bar` databases.

    * show query form for `db_pq`:

        ```
        /db://db_pq/
        ```

* `db-dump` - dump tables to XLSX.

    ```console
    curl http://127.0.0.1:9280/db-dump://mysql/ -o dst.xlsx
    ```

    Or access `http://127.0.0.1:9280/db-dump://mysql/` by web browser.
    Then start to download a excel file.

* `db-restore` - restore (clear all and import) tables from XLSX.

    ```console
    curl http://127.0.0.1:9280/db-restore://mysql/ -X POST \
      -H 'Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' \
      --data-binary @src.xlsx
    ```

    Or access `http://127.0.0.1:9280/db-dump://mysql/` by web browser.
    You can upload a excel file from the form.

* `db-update` - update tables by XLSX (upsert)

    ```console
    curl http://127.0.0.1:9280/db-update://mysql/ -X POST \
      -H 'Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' \
      --data-binary @src.xlsx
    ```

    Or access `http://127.0.0.1:9280/db-dump://mysql/` by web browser.
    You can upload a excel file from the form.

* `redis` - access to redis.

    [See document for details](doc/protocol-redis.md).

* `trdsql` - TRDSQL query editor

    [See document for detail](doc/filter-trdsql.md).

* `echarts` - ECharts query editor

    [See document for detail](doc/filter-echarts.md).

* `duckdb` - DockDB WASM Shell

    [See docuemnt for detail](doc/protocol-duckdb.md).

* `opfs` - Operate with OPFS (Origin Private File System)

    [See document for detail](doc/protocol-opfs.md).

* `examples` - Example files to use demo/document of filters

* `config` - current nvgd's configuration

    `/config://` or `/config/` (alias)

* `help` - show help (README.md) of nvgd.

    `/help://` or `/help/` (alias)

    It would be better that combining with `markdown` filter.

    ```
    /help/?markdown
    ```

* `version` - show nvgd's version

    Path is `/version://` or `/version/` (alias)

See also:

* [Filters](#filters)


## Configuration file

Nvgd takes a configuration file in YAML format.  A file `nvgd.conf.yml` in
current directory or given file with `-c` option is loaded at start.

`nvgd.conf.yml` consist from these parts:

```yml
# Listen IP address and port (OPTIONAL, default is "127.0.0.1:9280")
addr: "0.0.0.0:8080"

# Enable TLS (HTTPS) to serve. See doc/secure-contexts.md for details.
tls: (...sniped...)
    
# Path prefix for absolute links, use for sub-path multiple tenancy
path_prefix: /tenant_name/

# Destination (path or keyword) for error log, default is `(stderr)`
error_log: (stderr)

# Destination (path or keyword) for access log, default is `(discard)`
access_log: (stdout)

# File which served as "/" root.
root_contents_file: "/opt/nvgd/index.html"

# Configuration for protocols (OPTIONAL)
protocols:

  # File protocol handler's configuration.
  file:
    ...

  # Pre-defined command handlers.
  command:
    ...

  # AWS S3 protocol handler configuration (see other section, OPTIONAL).
  s3:
    ...

  # DB protocol handler configuration (OPTIONAL, see below)
  db:
    ...

# Configuration for each filters (OPTIONAL)

  indexhtml:
    ...

  htmltable:
    ...

  markdown:
    ...

# Default filters: pair of path prefix and filter description.
default_filters:
  ...

# Custom prefix aliases, see later "Prefix Aliases" section.
aliases:
  ...

# Enable "Cross-Origin Resource Sharing" (CORS).
# This value is put with "Access-Control-Allow-Origin" header in responses.
access_control_allow_origin: "*"
```

### File Protocol Handler

Example:

```yaml
file:
  locations:
    - '/var/log/'
    - '/etc/'

  forbiddens:
    - '/etc/ssh'
    - '/etc/passwd'

  use_unixtime: true
```

This configuration has `locations`, `forbiddens` properties.  These props
define accessible area of file system.

When paths are given as `locations`, only those paths are permitted to access,
others are forbidden.  Otherwise, all paths are accessible.

When `forbiddens` are given, those paths can't be accessed even if it is under
path in `locations`.

If the value of `use_unixtime` property is set to true, UNIX time will be used
instead of RFC1123 for all time expressions: `modified_at` or so.

For TLS (HTTPS) enablement, please refer to the documentation: [Secure Contexts](doc/secure-contexts.md).

### Command Protocol Handlers

Configuration of pre-defined command protocol handler maps a key to
corresponding command source.

Example:

```yml
command:
  "df": "df -h"
  "lstmp": "ls -l /tmp"
```

This enables two resources under `command` protocol.

* `/command://df`
* `/command://lstmp`

You could add filters of course, like: `/command://df?grep=re:foo`

### Config S3 Protocol Handlers

Configuration of S3 protocor handlers consist from 2 parts: `default` and
`buckets`.  `default` part cotains default configuration to connect S3.  And
`buckets` part could contain configuration for each buckets specific.

```yml
s3:

  # IANA timezone to show times (optional).  "Asia/Tokyo" for JST.
  timezone: Asia/Tokyo

  # default configuration to connect to S3 (REQUIRED for S3)
  default:

    # Access key ID for S3 (REQUIRED)
    access_key_id: xxxxxxxxxxxxxxxxxxxx

    # Secret access key (REQUIRED)
    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

    # Access point to connect (OPTIONAL, default is "ap-northeast-1")
    region: ap-northeast-1

    # Session token to connect (OPTIONAL, default is empty: not used)
    session_token: xxxxxxx

    # MaxKeys for S3 object listing. valid between 1 to 1000.
    # (OPTIONAL, default is 1000)
    max_keys: 10

    # HTTP PROXY to access S3. (OPTIONAL, default is empty: direct access)
    http_proxy: "http://your.proxy:port"

  # bucket specific configurations (OPTIONAL)
  buckets:

    # bucket name can be specified as key.
    "your_bucket_name1":
      # same properties with "default" can be placed at here.

    # other buckets can be added here.
    "your_bucket_name2":
      ...

  # UNIX time will be used instead of RFC1123 for all time expression:
  # `modified_at` or so. (OPTIONAL)
  use_unixtime: true
```

### Config DB Protocol Handler

Sample of configuration for DB protocol handler.

```yml
db:
  # key could be set favorite name for your database
  db_pq:
    # driver supports 'postgres' or 'mysql' for now
    driver: 'postgres'
    # name is driver-specific source name (DSN)
    name: 'postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'
    # limit number of rows for a query (default: 100)
    max_rows: 50

  # sample of connecting to MySQL
  db_mysql:
    driver: 'mysql'
    name:   'user:password@/dbname'
```

With above configuration, you will be able to access those databases with below URLs or commands.

  * `curl 'http://127.0.0.1:9280/db://db_pq/select%20email%20from%20users'`
  * `curl 'http://127.0.0.1:9280/db://db_mysql/select%20email%20from%20users'`

#### MySQL: TRADITIONAL mode

To restore or update MySQL database with `db-restore` or `db-update` protocol,
we recommend to use TRADITIONAL mode to make MySQL checks types strictly.  You
should add `?sql_mode=TRADITIONAL` to connection URL to enabling it.

Example:

```yaml
  db:
    mysql1:
      driver: mysql
      name: "mysql:abcd1234@tcp(127.0.0.1:3306)/mysql?sql_mode=TRADITIONAL"
```

#### Multiple Databases in an instance

To make DB protocol handler connect with multiple databases in an instance,
there are 3 steps to make it enable.

1.  Add `multiple_database: true` property to DB configuration.
2.  Add `{{.dbname}}` placeholder in value of `name`.
3.  Access to URL `/db://DBNAME@db_pq/you query`.

    DBNAME is used to expand `{{.dbname}}` in above.

As a result, your configuration would be like this:

```yml
db:
  db_pq:
    driver: 'postgres'
    name: 'postgres://pqgotest:password@localhost/{{.dbname}}?sslmode=verify-full'
    multiple_database: true

  # sample of connecting to MySQL
  db_mysql:
    driver: 'mysql'
    name:   'user:password@/{{.dbname}}'
    multiple_database: true
```

### Configure filters

Some filters can be configured by `filters` section.

* `custom_css_urls`: Array of string. Specify external CSS's URL for each string.

    Supported filters: htmltable, indexhtml, markdown

    Example: markdown filter outputs two `link` elements to including external CSS.
    indexhtml filters outputs a `link` elements for CSS.

    ```yaml
    filters:
      markdown:
        custom_css_urls:
        - https://www.kaoriya.net/assets/css/contents.css
        - https://www.kaoriya.net/assets/css/syntax.css

      indexhtml:
        custom_css_urls:
        - https://www.kaoriya.net/assets/css/contents.css
    ```

### Default Filters

Default filters provide a capability to apply implicit filters depending on
path prefixes. See [Filters](#filters) for detail of filters.

To apply `tail` filter for under `/file:///var/log/` path:

```yaml
default_filters:
  "file:///var/log/":
    - "tail"
```

If you want to show last 100 lines, change like this:

```yaml
default_filters:
  "file:///var/log/":
    - "tail=limit:100"
```

You can specify different filters for paths.

```yaml
default_filters:
  "file:///var/log/":
    - "tail"
  "file:///tmp/":
    - "head"
```

Default filters can be ignored separately by [all (pseudo) filter](#all-pseudo-filter).

Default filters are ignored for directories source of file protocols.


## Filters

Nvgd supports these filters:

* [Grep filter](#grep-filter)
* [Head filter](#head-filter)
* [Tail filter](#tail-filter)
* [Cut filter](#cut-filter)
* [Cutline filter](#cutline-filter)
* [Pager filter](#pager-filter)
* [Hash filter](#hash-filter)
* [LTSV Extraction filter](#ltsv-extraction-filter)
* [LTSV Conversion filter](#ltsv-conversion-filter)
* [JSONArray filter](#jsonarray-filter)
* [Index HTML filter](#index-html-filter)
* [Download to OPFS filter](#download-to-opfs-filter)
* [HTML Table filter](#html-table-filter)
* [Text Table filter](#text-table-filter)
* [Markdown filter](#markdown-filter)
* [Syntax Highlight filter](#syntax-highlight-filter)
* [TRDSQL filter](#trdsql-filter)
* [Echarts filter](#echarts-filter)
* [Refresh (pseudo) filter](#refresh-pseudo-filter)
* [Download (pseudo) filter](#download-pseudo-filter)
* [All (pseudo) filter](#all-pseudo-filter)

### Filter Spec

Where `{filters}` is:

    {filter}[&{filter}...]

Where `{filter}` is:

    {filter_name}[={options}]

Where `{options}` is:

    {option_name}:{value}[;{option_name}:{value}...]

See other section for detail of each filters.

Example: get last 50 lines except empty lines.

    /file:///var/log/messages?grep=re:^$;match:false&tail=limit:50

### Grep filter

Output lines which matches against regular expression.

As default, matching is made for whole line.  But when valid option `field` is
given, then matching is made for specified a field, which is splitted by
`delim` character.

`grep` command equivalent.

* filter\_name: `grep`
* options
    * `re` - regular expression used for match.
    * `match` - output when match or not match.  default is true.
    * `field` - a match target N'th field counted from 1.
        default is none (whole line).
    * `delim` - field delimiter string (default: TAB character).
    * `context` - show a few lines before and after the matched line.
        default is `0` (no contexts).
    * `number` - prefix each line of output with the 1-based line number.
        when `true`

### Head filter

Output the first N lines.

`head` command equivalent.

* filter\_name: `head`
* options
    * `start` - start line number for output.  default is 0.
    * `limit` - line number for output.  default is 10.

### Tail filter

Output the last N lines.

`tail` command equivalent.

* filter\_name: `tail`
* options
    * `limit` - line number for output.  default is 10.

### Cut filter

Output selected fields of lines.

`cut` command equivalent.

* filter\_name: `cut`
* options:
    * `delim` - field delimiter string (default: TAB character).
    * `white` - use consecutive whites as one single field separator (default: false)
    * `list` - selected fields, combinable by comma `,`.
        * `N` - N'th field counted from 1.
        * `N-M` - from N'th, to M'th field (included).
        * `N-` - from N'th field, to end of line.
        * `-N` - from first, to N'th field.

### Cutline filter

Cutlines extracts specific range of lines with regular expressions.
[See document for details](doc/filter-cutline.md).

### Pager filter

`pager` is a filter that divides the input stream into pages by lines that
match the specified pattern.

* filter\_name: `pager`
* options:
    * `eop`: Regular expression that matches page separator lines.
    * `pages`: Page number to output (1-based number)

        You can specify multiple pages separated by commas. Examples

        * `1`: First page only
        * `2,4,6`: Page 2, 4, and 6
        * `-1`: Last page
        * `-3`: 3rd page from the end
        * `1,-1`: First and last pages
        * `10-12`: Pages 10 to 12

    * `num`: Boolean. Output a page number at the top of the page.

      Example: `(page 12)`

### Hash filter

Output hash value.

* filter\_name: `hash`
* options:
    * `algorithm` - one of `md5` (default), `sha1`, `sha256` or `sha512`
    * `encoding` - one of `hex` (default), `base64` or `binary`

### Count filter

Count lines.

* filter\_name: `count`
* options: (none)

### LTSV Extraction filter

A filter that outputs only the specified labels from the rows of LTSV that
match the another specified label value.

* filter\_name: `ltsv`
* options:
    * `grep` - match parameter: `{label},{pattern}`
    * `match` - output when match or not match.  default is true.
    * `cut` - selected labels, combinable by comma `,`.

### LTSV Conversion filter

A filter to convert LTSV to another format. You can specify the output format
with the `format` parameter. The default is `tsv`.

* filter\_name: `ltsvconv`
* options:
    * `format` - output format (`tsv`, `csv`. default: `tsv`)

### JSONArray filter

Convert each line as a string of JSON array.

* filter\_name: `jsonarray`
* options: (none)

### Index HTML filter

Convert LTSV to Index HTML.
(limited for s3list and files (dir) source for now)

* filter\_name: `indexhtml`
* options:
    * `timefmt`: Time layout for "Modified At" or so. default is `RFC1123`.
      Possible values are, case insensitive: `ANSIC`, `UNIX`, `RUBY`, `RFC822`,
      `RFC822Z`, `RFC850`, `RFC1123`, `RFC1123Z`, `RFC3339`, `RFC3339NANO`,
      `STAMP`, and `DATETIME`
    * `nouplink`: Hide the "Up" link to navigate back through the directory
      hierarchy, when its value is `true`. Default is `false`: show the "Up"
      link.
    * `noopfs`: Hide links to OPFS uploader, when its value is `true`. Default
      is `false` show links to OPFS uploader.
* configurations:
    * `custom_css_urls`: list of URLs to link as CSS.

Example: list objects in S3 bucket "foo" with Index HTML.

    http://127.0.0.1:9280/s3list://foo/?indexhtml

This filter should be the last of filters.

### Download to OPFS filter
Provides a UI for downloading files from the LTSV file list to OPFS.
(limited for s3list and files (dir) source for now)

* filter\_name: `toopfs`
* optoins: (none)

Example: UI for download objects in S3 bucket "foo".

    http://127.0.0.1:9280/s3list://foo/?toopfs

### HTML Table filter

Convert LTSV to HTML table.

* filter\_name: `htmltable`
* options:
    * `linefeed` - boolean: expand all `\n` as linefeed.
* configurations:
    * `custom_css_urls`: list of URLs to link as CSS.

Example: query id and email column from users table on mine database.

    http://127.0.0.1:9280/db://mine/select%20id,email%20from%20users?htmltable

This filter should be the last of filters.

### Text Table filter

Convert LTSV to plain text table.

* filter\_name: `texttable`
* options: (none)

Example: query id and email column from users table on mine database.

    http://127.0.0.1:9280/db://mine/select%20id,email%20from%20users?texttable

Above query generate this table.

    +-----+-----------------------+
    |  id |        email          |
    +-----+-----------------------+
    |    0|foo@example.com        |
    |    1|bar@example.com        |
    +-----+-----------------------+

This filter should be the last of filters.

### Markdown filter

Convert markdown text to HTML.

* filter\_name: `markdown`
* options: (none)
* configurations:
    * `custom_css_urls`: list of URLs to link as CSS.

Example: show help in HTML.

    http://127.0.0.1:9280/help://?markdown
    http://127.0.0.1:9280/help/?markdown

### Syntax Highlight filter

Syntax Highlighting content.

*   filter\_name: `highlight`
*   options:
    *   `lexer`: Specifying the syntax lexer. A list of all lexers available can be found in [pygments-lexers.txt](https://github.com/alecthomas/chroma/blob/master/pygments-lexers.txt)
    *   `style`: Specify the style to display. For a quick overview of the available styles and how they look, check out [the Chroma Style Gallery](https://xyproto.github.io/splash/docs/all.html).

Example: syntax highlighting for Markdown instructions and HTML conversions.

    http://127.0.0.1:3000/help://?highlight=lexer:markdown
    http://127.0.0.1:3000/help://?markdown&highlight=lexer:html

### TRDSQL filter

TRDSQL filter provides SQL on CSV.
[See document for detail](doc/filter-trdsql.md).

### Echarts filter

Echarts filter provides drawing charts feature.
[See document for detail](doc/filter-echarts.md).

### Refresh (pseudo) filter

Add "Refresh" header with specified time (sec).

* filter\_name: `refresh`
* options: interval seconds to refresh.  0 for disable.

Example: Open below URL using WEB browser, it refresh in each 5 seconds
automatically.

    http://127.0.0.1:9280/file:///var/log/messages?tail&refresh=5

### Download (pseudo) filter

Add "Content-Disposition: attachment" header to the response.  It make the
browser to download the resource instead of showing in it.

* filter\_name: `download`
* options: (none)

Example: download the file "messages" and would be saved as file.

    http://127.0.0.1:9280/file:///var/log/messages?download

### All (pseudo) filter

Ignore [default filters](#default-filters)

* filter\_name: `all`
* options: (none)

Example: if specified some default filters for `file:///var/`, this ignore
those.

    http://127.0.0.1:9280/file:///var/log/messages?all


## Prefix Aliases

nvgd supports prefix aliases to keep compatibilities with [koron/night][night].
Currently these aliases are registered.

* `files/` -> `file:///`
* `commands/` -> `command://`
* `config/` -> `config://`
* `help/` -> `help://`
* `trdsql/` -> `trdsql:///`
* `echarts/` -> `echats:///`
* `version/` -> `version://`

For example this URL:

    http://127.0.0.1:9280/files/var/log/messages

It works same as below URL:

    http://127.0.0.1:9280/file:///var/log/messages

### Custom prefix aliases

You can add custom prefix aliases with `aliases` section in `nvgd.conf.yml`.

For example with below settings...

```yaml
aliases:
  'dump/': 'db-dump://'
```

You can dump a "mytable" table in "mydb" RDBMS with this URL:

    http://127.0.0.1:9280/dump/mydb/mytable

Instead of this:

    http://127.0.0.1:9280/db-dump://mydb/mytable

Custom prefix aliases can be used to avoid to containing `://` sub string in
path.


## References

* [koron/night][night] previous implementation which written in NodeJS.

[night]:https://github.com/koron/night
[globspec]:https://golang.org/pkg/path/filepath/#Match
