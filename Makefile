SHELL := /bin/bash

.PHONY: test test-local build

all: test build local-run

test:
	go test ./src/**/

build:
	docker build --tag new-books-notification .

deploy:
	. ./.env_prod && gcloud builds submit --config=cloudbuild.yaml \
  --substitutions=_GCP_BIGQUERY_DATASET="$${GCP_BIGQUERY_DATASET}",_GCP_BIGQUERY_TABLE="$${GCP_BIGQUERY_TABLE}",_GCS_BUCKET_NAME="$${GCS_BUCKET_NAME}",_SECRETS_NAME="$${SECRETS_NAME}" .

set-secrets:
	# If you set secrets for the first time, use "gcloud secrets create" instead of "gcloud secrets versions add"
	. ./.env_prod && echo -n "$${SLACK_WEBHOOK_URL}" | gcloud secrets versions add $${SECRETS_NAME} \
    --data-file=-

local-run:
	docker run --env-file .env new-books-notification