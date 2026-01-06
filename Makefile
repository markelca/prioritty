build:
	go build -o bin/pt cmd/*.go

test:
	go test ./... -v

