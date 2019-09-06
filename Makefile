
.PHONY: build
build:
	go build -o build/many

.PHONY: run
run: build
	./build/many

.PHONY: test
test:
	go test
