build:
	go build -o development-trains ./cmd/main.go
test:
	./development-trains -i ./tests/sample.txt --summary
test1:
	./development-trains ./tests/test1.txt
dev:
	make build && make test