SHELL := /bin/bash

.PHONY: test test-local build

all: test build local-run

test:
	go test ./src/**/

build:
	docker build --tag new-books-notification .

local-run:
	docker run --env-file .env new-books-notification