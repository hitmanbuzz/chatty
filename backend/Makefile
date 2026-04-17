BINARY_NAME=chat-app
MAIN_PKG_PATH=./cmd/main.go

.PHONY: build
build:
	go build -o build/${BINARY_NAME} ${MAIN_PKG_PATH}

.PHONY: run
run:
	./build/${BINARY_NAME}

.PHONY: clean
clean:
	rm -rf build/
	go clean
