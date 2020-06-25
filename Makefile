version = 0.1.0
run_opts ?=
binary_prefix = gvm-

.PHONY: test
test: app
	go test -v -covermode=count -coverprofile=coverage.out ./...

.PHONY: viewcoverage
viewcoverage: coverage.out
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	go clean
	rm ./$(binary_prefix)* || true

.PHONY: app
app: clean
	gofmt -w -s .
	go vet
	go build -o $(binary_prefix)$(version)

.PHONY: run
run: gvm-$(version)
	./gvm-$(version) $(run_opts)
