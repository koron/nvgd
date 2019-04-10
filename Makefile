GO_SUBPKGS = $(shell go list ./... | grep -v /vendor/ | sed -e "s!$$(go list)!.!")
GOOS = $(shell go env GOOS)
GOARCH = $(shell go env GOARCH)
BINARCH_BASE = nvgd1.0.0.$(GOOS)_$(GOARCH)
BINARCH_OUTDIR = dist

default: build

build:
	go build -v .

test:
	go test $(GO_SUBPKGS)

test-full:
	go test -v -race $(GO_SUBPKGS)

lint:
	@echo "go vet"
	@go vet $(GO_SUBPKGS)
	@echo ""
	@echo "golint"
	@for f in $(GO_SUBPKGS) ; do golint $$f ; done
	@echo ""

cyclo:
	-gocyclo -top 10 -avg $(GO_SUBPKGS)
	@echo ""

cyclo-report:
	@echo gocyclo -over 14 -avg
	-@gocyclo -over 14 -avg $(GO_SUBPKGS)
	@echo ""

misspell:
	@echo misspell
	@find $(GO_SUBPKGS) -maxdepth 1 -type f | xargs misspell
	@echo ""

report: misspell cyclo-report lint

list-packages:
	@echo $(GO_SUBPKGS)

deps:
	go get -v -u -d -t ./...

binarch: build
	rm -rf $(BINARCH_OUTDIR)/$(BINARCH_BASE)
	mkdir -p $(BINARCH_OUTDIR)/$(BINARCH_BASE)
	cp nvgd README.md $(BINARCH_OUTDIR)/$(BINARCH_BASE)
	tar czfC $(BINARCH_OUTDIR)/$(BINARCH_BASE).tar.gz $(BINARCH_OUTDIR) $(BINARCH_BASE)
.PHONY: bin_archive

tags:
	gotags -f tags -R .
.PHONY: tags

.PHONY: test test-full lint cyclo report deps
