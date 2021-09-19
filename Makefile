VERSION_INFO=Version $(shell git describe --tags `git rev-list --tags --max-count=1`) with commit $(shell git describe --always --tags) built at $(shell date) on $(shell hostname) with $(shell go version)
GOPKGS =  ./internal/... ./cmd/...
GODIRS = internal cmd
MAINPKG = ./cmd/proxy

.PHONY: proxy proxy+cover fmt check check-fmt vet test ssl

proxy:
	go build -ldflags="-X \"main.versionInfo=$(VERSION_INFO)\"" -o ./bin/proxy $(MAINPKG)

ssl:
	/bin/zsh scripts/certificatesAndKeys.sh

proxy+cover:
	go test -c -o bin/$@ -tags testrunmain $(MAINPKG) -coverpkg ./pkg/... ./cmd/...

fmt:
	gofmt -l -w $(GODIRS)

vet:
	go vet $(GOPKGS)

check-fmt:
	test -z "`gofmt -l $(GODIRS)`"

check: check-fmt check-imports test vet

imports:
	goimports -l -w $(GODIRS)

check-imports:
	test -z "`goimports -l $(GODIRS)`"

test:
	mkdir -p coverage
	go test -race -coverprofile=coverage/unit.cover.out $(GOPKGS)

show_coverage:
	go tool cover -html=coverage/unit.cover.out -o coverage/cover.html

/DEFAULT_GOAL := proxy