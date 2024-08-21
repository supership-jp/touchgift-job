include ./localstack/Makefile
include ./dynamodb/Makefile
include ./sns/Makefile

.DEFAULT_GOAL := help
SHELL=/bin/bash

APPLICATION = touchgift-delivery-job
ECS_CLUSTER = touchgift-job
SERVICE = $(APPLICATION)

AWS_PROFILE = dummy
AWS_REGION=ap-northeast-1

# Dynamo設定まわり
DYNAMODB_ENDPOINT = http://localhost:8000
DYNAMODB_ENDPOINT_OPTS = --endpoint-url $(DYNAMODB_ENDPOINT)
DYNAMODB_OPTIONS = $(DYNAMODB_ENDPOINT_OPTS) --profile $(AWS_PROFILE)

# TODO: アプリケーション名, ecsクラスター, サービス名を設定
APPLICATION =
ECS_CLUSTER =
SERVICE =

CURRENT_TAG = 1.0.0
TAG = $(CURRENT_TAG)
ifndef TAG
override TAG = $(CURRENT_TAG)
endif

TEST_TARGET=./...

# TODO: リポジトリ名を設定
REPOSITORY = touchgift-job/$(APPLICATION)
PREFIX = $(REPOSITORY)
REGISTRY =

CREDENTIAL = tpo-$(ENV)
AWS_VAULT_CMD = aws-vault exec $(CREDENTIAL) --

# DynamoDB環境変数
DYNAMODB_ENDPOINT = http://localhost:4566
DYNAMODB_ENDPOINT_OPTS = --endpoint-url $(DYNAMODB_ENDPOINT)
DYNAMODB_OPTIONS = $(DYNAMODB_ENDPOINT_OPTS) --profile $(AWS_PROFILE)

#　SQS環境変数
SQS_ENDPOINT = http://localhost:4566
SQS_ENDPOINT_OPTS = --endpoint-url $(SQS_ENDPOINT)
SQS_OPTIONS = $(SQS_ENDPOINT_OPTS) --profile $(AWS_PROFILE)

# SNS環境変数
SNS_ENDPOINT = http://localhost:4566
SNS_ENDPOINT_OPTS = --endpoint-url $(SNS_ENDPOINT)
SNS_OPTIONS = $(SNS_ENDPOINT_OPTS) --profile $(AWS_PROFILE)

## 本番アカウント
### e(本番), rc(リリース前テスト環境)
## 検証アカウント
### e(検証), beta(ベータ環境)
ENVIRONMENT_SUFFIX = e

# suffixで環境変数を切り替える
ifneq (,$(findstring beta,$(ENVIRONMENT_SUFFIX)))
	DELIVERY_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=touchgift_creative_data.$(ENVIRONMENT_SUFFIX)
	DELIVERY_BUDGET_TABLE_NAME=touchgift_delivery_budget_data.$(ENVIRONMENT_SUFFIX)
	TOUCHGIFT_DB_HOST=/rds/touchgift-provisioned-beta/endpoint
	TOUCHGIFT_DB_PASSWORD=/rds/touchgift-provisioned-beta/user/touchgift/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-delivery-control-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/touchgift-delivery-operation-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/touchgift-delivery-control-$(ENVIRONMENT_SUFFIX)
else
	DELIVERY_TABLE_NAME=touchgift_delivery_data
	CREATIVE_TABLE_NAME=touchgift_creative_data
	DELIVERY_BUDGET_TABLE_NAME=touchgift_delivery_budget_data
	TOUCHGIFT_DB_HOST=/rds/touchgift-provisioned/endpoint
	TOUCHGIFT_DB_PASSWORD=/rds/touchgift-provisioned/user/touchgift/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-delivery-control
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/touchgift-delivery-operation
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/touchgift-delivery-control
endif
ifneq (,$(findstring rc,$(ENVIRONMENT_SUFFIX)))
	DELIVERY_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=touchgift_creative_data.$(ENVIRONMENT_SUFFIX)
	DELIVERY_BUDGET_TABLE_NAME=touchgift_delivery_budget_data.$(ENVIRONMENT_SUFFIX)
	TOUCHGIFT_DB_HOST=/rds/touchgift-provisioned-rc/endpoint
	TOUCHGIFT_DB_PASSWORD=/rds/touchgift-provisioned-rc/user/touchgift/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-delivery-control-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/touchgift-delivery-operation-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/touchgift-delivery-control-$(ENVIRONMENT_SUFFIX)
endif

build-image: aws_account_id ## Build docker image ex. make build-image ENV=staging TAG=
	@docker build --tag $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG) .

push-image: aws_account_id ## Push docker image ex. make push-image ENV=staging TAG=
	@$(AWS_VAULT_CMD) aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com && \
	docker push $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG) && \
	docker logout https://$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

output-image: aws_account_id ## output full image name make output-image ENV=staging
	@echo $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)

run: ## Start dev server
	@go build -o main && \
		GIN_MODE=release \
		LOG_LEVEL=info \
		SQS_ENDPOINT=$(SQS_ENDPOINT) \
		DYNAMODB_ENDPOINT=$(DYNAMODB_ENDPOINT) \
		TABLE_NAME_SUFFIX= \
		AWS_PROFILE=$(AWS_PROFILE) \
		./main || exit 0

tests: ## Test (キャッシュしたくない場合: clean-testcacheを実行)
	@AWS_PROFILE=$(AWS_PROFILE) \
		GIN_MODE=test \
		TABLE_NAME_PREFIX=$(TABLE_NAME_PREFIX) \
		DYNAMODB_ENDPOINT=$(DYNAMODB_ENDPOINT) \
		go test -cover -v $(TEST_TARGET)

clean-testcache: ## clean testcache
	go clean -testcache

generate:
	@go generate ./...

mod-tidy: ## packageのinstall
	go mod tidy -v

aws_account_id: ## current aws_account_id
	$(eval AWS_ACCOUNT_ID := $(shell $(AWS_VAULT_CMD) aws sts get-caller-identity --query 'Account' --output text))

help: ## Help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort |\
		awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## lint
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.54.2 golangci-lint run -v --timeout=5m

start-db: ## start db (mysql)
	@cd ./rdb && docker run --rm -it \
		-p 3306:3306 \
		-e MYSQL_ROOT_PASSWORD=test \
		-e MYSQL_DATABASE=touchgift \
		-e MYSQL_USER=touchgift \
		-e MYSQL_PASSWORD=test \
		-e TZ=Asia/Tokyo \
		-v $$(PWD)/init:/docker-entrypoint-initdb.d \
		--name touchgift-job-db \
		mysql:8.0
stop-db:
	docker stop touchgift-job-db

cleanup-db:
	make stop-db && \
	make start-db

create-db-testdata: ## create test data for rdb ex. make create-db-testdata
	AWS_PROFILE=$(AWS_PROFILE) go test -v ./infra -run TestData -tags=createdata -count 1
