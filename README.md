# NVGD - Night Vision Goggles Daemon

HTTP file server to help DevOps.

Index:

  * [Acceptable path](#acceptable-path)
  * [Configuration file](#configuration-file)
  * [Filters](#filters)

## How to use

Install:

    $ go get github.com/koron/nvgd

Run:

    $ nvgd

Access:

    $ curl http://127.0.0.1:9280/file:///var/log/message/httpd.log?tail=limit:25

Update:

    $ go get -u github.com/koron/nvgd


## Acceptable path

Nvgd accepts path in these like format:

    /{protocol}://{args/for/protocol}[?{filters}]

Nvgd supports these `protocol`s:

  * `file` - `/file:///path/to/source`
  * `s3` - `/s3://bucket-name/key/to/object`

Supported filters are:

  * [Grep filter](#grep-filter)
  * [Head filter](#head-filter)
  * [Tail filter](#tail-filter)
  * [Cut filter](#cut-filter)
  * [Hash filter](#hash-filter)
  * [LTSV filter](#ltsv-filter)
  * [Refresh filter](#refresh-filter)


## Configuration file

Nvgd takes a configuration file in YAML format.  A file `nvgd.conf.yml` in
current directory or given file with `-c` option is loaded at start.

`nvgd.conf.yml` consist from these parts:

```yml
# Listen IP address and port (OPTIONAL, default is "127.0.0.1:9280")
addr: "0.0.0.0:8080"

# Configuratio for protocols (OPTIONAL)
protocols:

  # AWS S3 protocol handler configuration (see other section, OPTIONAL).
  s3:
    ...
```

### Config S3 Protocol Handler

Configuration of S3 protocol handler consist from 2 parts: `default` and
`buckets`.  `default` part cotains default configuration to connect S3.  And
`buckets` part could contain configuration for each buckets specific.

```yml
s3:

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

  # bucket specific configurations (OPTIONAL)
  buckets:

    # bucket name can be specified as key.
    "your_bucket_name1":
      # same properties with "default" can be placed at here.

    # other buckets can be added here.
    "your_bucket_name2":
      ...
```


## Filters

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
    * `start` - start line number for output.  begging 0.  default is 0.
    * `limit` - line number for output.  defualt is 10.

### Tail filter

Output the last N lines.

`tail` command equivalent.

  * filter\_name: `tail`
  * options
    * `limit` - line number for output.  defualt is 10.

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

### Refresh  filter

Add "Refresh" header with specified time (sec).

  * filter\_name: `refresh`
  * options: interval seconds to refresh.  0 for disable.

Example: Open below URL using WEB browser, it refresh in each 5 seconds
automatically.

    http://127.0.0.1:9280/file:///var/log/messages?tail&refresh=5
