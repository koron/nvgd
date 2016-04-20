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
