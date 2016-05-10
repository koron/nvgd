# NVGD - Night Vision Goggles Daemon

HTTP file server to help DevOps.

## How to use

Install:

    $ go get github.com/koron/nvgd

Run:

    $ nvgd

Access:

    $ curl http://127.0.0.1:9280/file:///var/log/message/httpd.log?tail=limit:25

Update:

    $ go get -u github.com/koron/nvgd

## Configuration file

nvgd takes a configuration file `nvgd.conf.yml` in current directory or given
with `-c` option.

nvgd.conf.yml should consist these parts:

```yml
# Listen IP address and port.
add: "0.0.0.0:8080"

# Configuratio for protocols.
protocols:

  # AWS S3 protocol handler configuration (see other section).
  s3:
    ...
```

### Config S3 Protocol Handler

Configuration of S3 protocol handler should consist from 2 parts: `default` and
`buckets`.  `default` part cotains default configuration to connect S3.  And
`buckets` part could contain configuration for each buckets specific.

```yml
s3:
  # default configuration to connect to S3
  default:

    # Access key ID for S3 (REQUIRED)
    access_key_id: xxxxxxxxxxxxxxxxxxxx

    # Secret access key (REQUIRED)
    secret_access_key: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

    # Access point to connect (OPTIONAL, default is "ap-northeast-1")
    region: ap-northeast-1

    # Session token to connect (OPTIONAL, default is empty: not used)
    session_token: xxxxxxx

  # bucket specific configuration.
  buckets:

    # bucket name can be specified as key.
    "your_bucket_name1":
      # same properties with "default" can be used at here.

    # other buckets can be added here.
    "your_bucket_name2":
      ...
```
