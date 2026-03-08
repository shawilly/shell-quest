start:
	go build ./cmd/shellquest/ && go run ./cmd/shellquest

test:
	go test -v ./...
