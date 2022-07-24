SHELL := /bin/bash

.PHONY: test-local build

all: build local-run

build:
	docker build --tag new-books-notification .

local-run:
	docker run --env-file .env new-books-notification