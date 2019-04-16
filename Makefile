VERSION ?= $(shell git describe --tags 2>/dev/null | cut -c 2-)
TEST_FLAGS ?=
REPO_OWNER ?= $(shell cd .. && basename "$$(pwd)")

GOPATH=$(CURDIR)/../../../../
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html:
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-add:
ifdef PACKAGE
	@$(GOPATHCMD) dep ensure -add -v $(PACKAGE)
else
	@echo "PACKAGE envvar is not defined"
endif

dep-update:
	@$(GOPATHCMD) dep ensure -v -update

dcup:
	docker-compose up -d

dcdn:
	docker-compose down --remove-orphans

vet:
	@$(GOPATHCMD) go vet ./...

lint:
	@$(GOPATHCMD) golint

fmt:
	@$(GOPATHCMD) go fmt ./...

.PHONY: build-cli clean test-short test test-with-flags deps html-coverage \
        restore-import-paths rewrite-import-paths list-external-deps release \
        docs kill-docs open-docs kill-orphaned-docker-containers dep-ensure \
        dep-update

SHELL = /bin/bash
RAND = $(shell echo $$RANDOM)
