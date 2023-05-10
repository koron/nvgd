# How to build on CentOS 7

## Build with native CentOS 7

Install gcc and git. gcc is required to build CGO enabled packages.
git is required to checkout nvgd's source code.  `sudo` command or root user's
privilege are required to install gcc and git.

```console
$ sudo yum install -y gcc git
```

Download Go and install.
See also <https://go.dev/doc/install>.

```console
$ curl -fSsLo go.linux-amd64.tar.gz https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
$ rm -rf /usr/local/go && tar -C /usr/local -xzf go.linux-amd64.tar.gz
$ export PATH=$PATH:/usr/local/go/bin
$ export PATH=$PATH:$(go env GOPATH)/bin
$ rm -f go.linux-amd64.tar.gz
```

Checkout nvgd source code and build it.  If you want to install a specific
version of nvgd, replace `-b main` with you prefered version.  (ex. `-b
v1.12.2`)

```console
$ git clone -b main --depth 1 https://github.com/koron/nvgd.git nvgd
$ cd nvgd
$ go install
```

`nvgd` executable file is installed at `~/go/bin`.  You can move/copy the file into
your prefered directory.

You may remove source code after `nvgd` executable is copied.

```console
$ cd ..
$ rm -rf nvgd
```
