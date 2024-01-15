# How to install nvgd as rc.d script for FreeBSD

## Install via binary package

1. Create a binary package

    ```console
    $ make package
    ```

2. Copy the package to a target machine (OPTIONAL)

    ```console
    $ scp nvgd-freebsd14-snapshot-bin.tar.bz2 foobar@machine:tmp
    ```

3. Install as root

    ```console
    $ sudo tar xvf nvgd-freebsd14-snapshot-bin.tar.bz2 -C /
    ```

4. Copy and edit configuration file

    ```console
    $ cd /usr/local/etc
    $ sudo cp -p nvgd.conf.yml{.sample,}    # OPTIONAL: copy sample as start
    $ sudo vi nvgd.conf.yml
    ```

    You should edit `addr` especially.

5. Add `nvgd_enable=YES` to /etc/rc.conf.local or /etc/rc.conf

    Example to add it:

    ```console
    $ sudo sysrc nvgd_enable=YES
    ```

6. Start nvgd servcie

    ```console
    $ sudo service nvgd start
    ```

## Install from source directly

1. Build nvgd and copy it to /usr/local/bin

    ```console
    $ go install github.com/koron/nvgd@latest
    $ sudo install -o root -g wheel -m 0555 ~/go/bin/nvgd /usr/local/bin/
    ```

2. Edit [nvgd.conf.yml](./nvgd.conf.yml) and copy it to /usr/local/etc

    You should edit `addr` especially.

    ```console
    $ sudo install -o root -g wheel -m 0644 nvgd.conf.yml /usr/local/etc/
    ```

3. Copy [rc.d-nvgd](./rc.d-nvgd) to /usr/local/etc/rc.d/nvgd

    ```console
    $ sudo install -o root -g wheel -m 0755 rc.d-nvgd /usr/local/etc/rc.d/nvgd
    ```

4. Copy [newsyslog.conf.d-nvgd.conf](./newsyslog.conf.d-nvgd.conf) to /usr/local/etc/newsyslog.conf.d/nvgd.conf

    ```console
    $ sudo install -o root -g wheel -m 0644 newsyslog.conf-nvgd.conf /usr/local/etc/newsyslog.conf.d/nvgd.conf
    ```

See next "Start and stop nvgd service" also.

## Start and stop nvgd service

1. Add `nvgd_enable="YES"` to /etc/rc.conf.local or /etc/rc.conf

2. Start with `sudo service nvgd start`

3. Stop with `sudo service nvgd stop`
