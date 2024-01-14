# How to install nvgd as rc.d script for FreeBSD

1. Build nvgd and copy it to /usr/local/bin

    ```console
    $ go install github.com/koron/nvgd@latest
    $ sudo install -o root -g wheel -m 0555 ~/go/bin/nvgd /usr/local/bin/
    ```

2. Edit [nvgd.conf.yml](./nvgd.conf.yml) and copy it to /usr/local/etc

    ```console
    $ sudo install -o root -g wheel -m 0644 nvgd.conf.yml /usr/local/etc/
    ```

2. Copy [rc.d-nvgd](./rc.d-nvgd) to /usr/local/etc/rc.d/nvgd

    ```console
    $ sudo install -o root -g wheel -m 0755 rc.d-nvgd /usr/local/etc/rc.d/nvgd
    ```

3. Copy [newsyslog.conf.d-nvgd.conf](./newsyslog.conf.d-nvgd.conf) to /usr/local/etc/newsyslog.conf.d/nvgd.conf

    ```console
    $ sudo install -o root -g wheel -m 0644 newsyslog.conf-nvgd.conf /usr/local/etc/newsyslog.conf.d/nvgd.conf
    ```
