# Redis protocol

## Gettings started

Supposed to redis is running on local.

Add this setting to `protocols` section in your `nvgd.conf.yml` then start
`nvgd`.

```yaml
  redis:
    stores:
      local:
        url: redis://127.0.0.1:6379/0
```

Set some strings to your redis for test.

```console
$ redis-cli set foo 123456789
OK

$ redis-cli set foo abcxyz
OK
```

Then try to get values from redis via nvgd using curl or so.

```console
$ curl http://127.0.0.1:9280/redis://local/get/foo
123456789

$ curl http://127.0.0.1:9280/redis://local/get/bar
abcxyz
```

## URL spec

    redis://{store_name}/{command}/[ARGUMENTS]

Where `store_name` should be replaced by one of names on
`protocols/reids/stores` section in `nvgd.conf.yml`.

Where `command` (case ignored):

*   `get` - get value(s) with key.
*   `keys` - list keys which match with ARGUMENTS as pattern.

## Get command

    redis://{store_name}/get/{key}[(/ARGUMENTS)*]

Behavior of get command will be changed by types for key. Supported types are:
`string`, `list`, `set`, `zset` and `hash`.

*   `string`
    *   0 arguments: like `GET {key}`

        Example: `redis://local/get/foo`

    *   1 argument: like `GETBIT {key} {offset}`

        Example: `redis://local/get/foo/8`

    *   2 arguments: like `GETRANGE {key} {start} {end}`

        Example: `redis://local/get/foo/1/3`

*   `list`
    *   0 arguments: like `LLEN {key}`
    *   1 argument: like `LINDEX {key} {index}`
    *   2 arguments: like `LRANGE {key} {start} {stop}`
*   `set`
    *   0 arguments: like `SCARD {key}`
    *   1 argument: like `SISMEMBER {key} {member}`
*   `zset`
    *   0 arguments: like `ZCARD {key}`
    *   1 argument: like `ZRANK {key} {member}`
    *   2 arguments: like `ZRANGE {key} {start} {stop}`
*   `hash`
    *   0 arguments: like `HLEN {key}`
    *   1 argument: like `HGET {key} {field}`

## Keys command

    redis://{store_name}/keys[/{PATTERN}]

Equivalent with `KEYS` command of redis.
<https://redis.io/commands/keys>

When the pattern doesn't include any meta characters, nvgd will append `*` at
last. It will help to implement type a head (incremental) search.

Examples:

* `redis://{store_name}/keys` - `KEYS *`
* `redis://{store_name}/keys/a` - `KEYS a*`
* `redis://{store_name}/keys/ab` - `KEYS ab*`
* `redis://{store_name}/keys/*a` - `KEYS *a`, no `*` supplied at tail.
* `redis://{store_name}/keys/*a*` - `KEYS *a*`
* `redis://{store_name}/keys/%3fa` - `KEYS ?a`, `?` should be URL encodeded.
