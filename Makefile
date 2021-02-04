lint:
	gofumpt -w -s ./..
	gofumports -w ./..
	golint ./...
	golangci-lint run --fix

test:
	go test ./...