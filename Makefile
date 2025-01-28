build:
	go build -o development-trains ./cmd/main.go

build-image:
	docker build -t idea456:development-trains .

test:
	./development-trains -i ./tests/sample.txt --summary

test-image:
	@docker run -v $(pwd)/tests:/tests idea456:development-trains -i ./tests/sample.txt --summary

dev:
	make build && make test