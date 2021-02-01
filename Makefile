lint:
	gofumpt -w -s ./..
	gofumports -w ./..
	golangci-lint run --fix

test:
	go test ./...