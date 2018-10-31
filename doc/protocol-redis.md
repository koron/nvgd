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

    redis://{store_name}/{command}/[ARGUMENS]

Where `store_name` should be replaced by one of names on
`protocols/reids/stores` section in `nvgd.conf.yml`.

Where `command` (case ignored):

*   `get` - get value(s) with key.

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
