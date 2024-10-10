build: 
	@go build -o bin/CommerceCore cmd/main.go

test: 
	@go test -v ./...

run:build
	@./bin/CommerceCore