.PHONY: test
test:
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: viewcoverage
viewcoverage: coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	go clean
	rm -rf package

.PHONY: build
build: clean
	gofmt -w -s .
	mkdir -p package
	cp -R bin package/bin
	go build -o package/bin/goenv-fetch cmd/fetch/main.go
	cp install-goenv.sh package

.PHONY: pkg
pkg: build
	cp -R ./package ./goenv
	tar -czf goenv.tar.gz ./goenv
	rm -rf ./goenv
