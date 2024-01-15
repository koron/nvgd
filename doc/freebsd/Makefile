VERSION = snapshot
DESTDIR = $(.CURDIR)/work
PREFIX	= /usr/local

OS_TARGET = freebsd14
PKGFILE   = nvgd-$(OS_TARGET)-$(VERSION)-bin.tar.bz2

.PHONY: package
package: clean test build install
	rm -f $(PKGFILE)
	tar caf $(PKGFILE) --uname root --gname wheel --strip-components 1 -C $(DESTDIR) .

.PHONY: info
info:
	@echo DESTDIR=$(DESTDIR)
	@echo VERSION=$(VERSION)

.PHONY: build
build:
	go build -v ../..

.PHONY: test
test:
	go test ../../...

.PHONY: clean
clean:
	rm -f nvgd
	rm -rf $(DESTDIR)

.PHONY: install
install:
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	mkdir -p $(DESTDIR)$(PREFIX)/etc/
	mkdir -p $(DESTDIR)$(PREFIX)/etc/rc.d
	mkdir -p $(DESTDIR)$(PREFIX)/etc/newsyslog.conf.d
	cp -p ./nvgd $(DESTDIR)$(PREFIX)/bin/nvgd
	chmod 0555 $(DESTDIR)$(PREFIX)/bin/nvgd
	cp -p nvgd.conf.yml $(DESTDIR)$(PREFIX)/etc/nvgd.conf.yml.sample
	chmod 0644 $(DESTDIR)$(PREFIX)/etc/nvgd.conf.yml.sample
	cp -p rc.d-nvgd $(DESTDIR)$(PREFIX)/etc/rc.d/nvgd
	chmod 0755 $(DESTDIR)$(PREFIX)/etc/rc.d/nvgd
	cp -p newsyslog.conf.d-nvgd.conf $(DESTDIR)$(PREFIX)/etc/newsyslog.conf.d/nvgd.conf
	chmod 0644 $(DESTDIR)$(PREFIX)/etc/newsyslog.conf.d/nvgd.conf
