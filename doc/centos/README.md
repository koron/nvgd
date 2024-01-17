# How to compile on CentOS

## Build instruction

1.  Install pre-requirements:

    ```sh
    # Install pre-required packages
    yum install -y centos-release-scl
    yum install -y gcc
    yum clean all

    # Install Go
    curl -fsSLO https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
    rm -f go1.21.6.linux-amd64.tar.gz 

    # Set PATH env to enable go command
    export PATH=$PATH:/usr/local/go/bin
    ```

2.  Build and install latest nvgd


    ```console
    $ go install github.com/koron/nvgd@latest
    ```

    This install nvgd executable into ~/go/bin dir.
