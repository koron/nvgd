# NVGD - Night Vision Goggles Daemon

HTTP file server to help DevOps.

[![Actions/Go](https://github.com/koron/nvgd/workflows/Go/badge.svg)](https://github.com/koron/nvgd/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/nvgd)](https://goreportcard.com/report/github.com/koron/nvgd)

Index:

  * [How to use](#how-to-use)
  * [Acceptable path](#acceptable-path)
    * [Protocols](#protocols)
  * [Configuration file](#configuration-file)
    * [File Protocol Handler](#file-protocol-handler)
    * [Command Protocol Handlers](#command-protocol-handlers)
    * [S3 Protocol Handlers](#config-s3-protocol-handlers)
    * [Config DB Protocol Handler](#config-db-protocol-handler)
    * [Default Filters](#default-filters)
  * [Filters](#filters)
  * [Prefix Aliases](#prefix-aliases)

## How to use

Install:

    $ go get github.com/koron/nvgd

Run:

    $ nvgd

Access:

    $ curl http://127.0.0.1:9280/file:///var/log/message/httpd.log?tail=limit:25

Update:

    $ go get -u github.com/koron/nvgd

### Extra install

GOPATH 外でビルドする場合は以下のようにする。
アップデートする場合も同様。
作成したバイナリは `$GOPATH/bin/nvgd` として保存される。

```console
$ GO111MODULE=on go get github.com/koron/nvgd
```

クロスコンパイルする場合は `GOOS` と `GOARCH` を指定する。
以下の例はLinux/ARM64 (例: ラズパイ4等) 向け。

```console
$ GOOS=linux GOARCH=arm64 GO111MODULE=on go get github.com/koron/nvgd
```

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

# Configuratio for protocols (OPTIONAL)
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

# Default filters: pair of path prefix and filter description.
default_filters:
  ...

# Custom prefix aliases, see later "Prefix Aliases" section.
aliases:
  ...
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
```

This configuration has `locations` and `forbiddens` properties.  These props
define accessible area of file system.

When paths are given as `locations`, only those paths are permitted to access,
others are forbidden.  Otherwise, all paths are accessible.

When `forbiddens` are given, those paths can't be accessed even if it is under
path in `locations`.

### Commnad Protocol Handlers

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
  * [Hash filter](#hash-filter)
  * [LTSV filter](#ltsv-filter)
  * [JSONArray filter](#jsonarray-filter)
  * [Index HTML filter](#index-html-filter)
  * [HTML Table filter](#html-table-filter)
  * [Text Table filter](#text-table-filter)
  * [Markdown filter](#markdown-filter)
  * [Refresh filter](#refresh-filter)
  * [Download filter](#download-filter)
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
    * `list` - selected fields, combinable by comma `,`.
      * `N` - N'th field counted from 1.
      * `N-M` - from N'th, to M'th field (included).
      * `N-` - from N'th field, to end of line.
      * `N-` - from first, to N'th field.

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

### LTSV filter

Output, match to value of specified label, and output selected labels.

  * filter\_name: `ltsv`
  * options:
    * `grep` - match parameter: `{label},{pattern}`
    * `match` - output when match or not match.  default is true.
    * `cut` - selected labels, combinable by comma `,`.

### JSONArray filter

Convert each line as a string of JSON array.

  * filter\_name: `jsonarray`
  * options: (none)

### Index HTML filter

Convert LTSV to Index HTML.
(limited for s3list and files (dir) source for now)

  * filter\_name: `indexhtml`
  * options: (none)

Example: list objects in S3 bucket "foo" with Index HTML.

    http://127.0.0.1:9280/s3list://foo/?indexhtml

This filter should be the last of filters.

### HTML Table filter

Convert LTSV to HTML table.

  * filter\_name: `htmltable`
  * options: (none)

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

Example: show help in HTML.

    http://127.0.0.1:9280/help://?markdown
    http://127.0.0.1:9280/help/?markdown

### Refresh filter

Add "Refresh" header with specified time (sec).

  * filter\_name: `refresh`
  * options: interval seconds to refresh.  0 for disable.

Example: Open below URL using WEB browser, it refresh in each 5 seconds
automatically.

    http://127.0.0.1:9280/file:///var/log/messages?tail&refresh=5

### Download filter

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
