SHELL := /bin/bash

.PHONY: test test-local build

all: test build local-run

test:
	go test ./src/**/

build:
	docker build --tag new-books-notification .

deploy:
	. ./.env_prod && gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_SLACK_WEBHOOK_URL="$${SLACK_WEBHOOK_URL}",_GCP_BIGQUERY_DATASET="$${GCP_BIGQUERY_DATASET}",_GCP_BIGQUERY_TABLE="$${GCP_BIGQUERY_TABLE}",_GCS_BUCKET_NAME="$${GCS_BUCKET_NAME}" .

local-run:
	docker run --env-file .env new-books-notification