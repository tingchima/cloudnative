server:
	GOFLAGS=-mod=vendor go run ./cmd server

gomod.update:
	go mod tidy
	go mod vendor