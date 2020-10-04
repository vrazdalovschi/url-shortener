CONFIG_PATH=${HOME}/.url-shortener/

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: test
test:
	go test -race ./...

TAG ?= 0.1.0

build-docker:
	docker build -t github.com/vrazdalovschi/url-shortener:$(TAG) .