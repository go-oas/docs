lint:
	gofumpt -w -s ./..
	golint ./...
	golangci-lint run --fix

test:
	go test ./...

update_cache:
	curl https://sum.golang.org/lookup/github.com/go-oas/docs@v$(VER)