build:
	go build -o development-trains ./cmd/main.go
test:
	./development-trains ./tests/sample.txt

test1:
	./development-trains ./tests/test1.txt