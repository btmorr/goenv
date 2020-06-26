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
	go build -o package/bin/gvm-fetch cmd/fetch/main.go
	go build -o package/bin/gvm-fetch-dry-run cmd/fetch-dry-run/main.go
	go build -o package/bin/gvm-version cmd/version/main.go
	cp install-gvm.sh package
