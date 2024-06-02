.DEFAULT_GOAL := help
SHELL=/bin/bash

# TODO: アプリケーション名, ecsクラスター, サービス名を設定
APPLICATION =
ECS_CLUSTER =
SERVICE =

CURRENT_TAG = 1.0.0
TAG = $(CURRENT_TAG)
ifndef TAG
override TAG = $(CURRENT_TAG)
endif

# TODO: リポジトリ名を設定
REPOSITORY =
PREFIX = $(REPOSITORY)
REGISTRY =

CREDENTIAL = tpo-$(ENV)
AWS_VAULT_CMD = aws-vault exec $(CREDENTIAL) --

AWS_PROFILE = dummy
AWS_REGION=ap-northeast-1

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
	DELIVERY_TABLE_NAME=dm_delivery_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=dm_creative_data.$(ENVIRONMENT_SUFFIX)
	DELIVERY_BUDGET_TABLE_NAME=dm_delivery_budget_data.$(ENVIRONMENT_SUFFIX)
	DM_DB_HOST=/rds/dm-provisioned-beta/endpoint
	DM_DB_PASSWORD=/rds/dm-provisioned-beta/user/dm/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):dm-delivery-control-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-operation-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-control-$(ENVIRONMENT_SUFFIX)
else
	DELIVERY_TABLE_NAME=dm_delivery_data
	CREATIVE_TABLE_NAME=dm_creative_data
	DELIVERY_BUDGET_TABLE_NAME=dm_delivery_budget_data
	DM_DB_HOST=/rds/dm-provisioned/endpoint
	DM_DB_PASSWORD=/rds/dm-provisioned/user/dm/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):dm-delivery-control
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-operation
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-control
endif
ifneq (,$(findstring rc,$(ENVIRONMENT_SUFFIX)))
	DELIVERY_TABLE_NAME=dm_delivery_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=dm_creative_data.$(ENVIRONMENT_SUFFIX)
	DELIVERY_BUDGET_TABLE_NAME=dm_delivery_budget_data.$(ENVIRONMENT_SUFFIX)
	DM_DB_HOST=/rds/dm-provisioned-rc/endpoint
	DM_DB_PASSWORD=/rds/dm-provisioned-rc/user/dm/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):dm-delivery-control-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-operation-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/adm-delivery-dm-delivery-control-$(ENVIRONMENT_SUFFIX)
endif

run: ## Start dev server
	@go build -o main && \
		GIN_MODE=release \
		LOG_LEVEL=info \
		SQS_ENDPOINT=$(SQS_ENDPOINT) \
		DYNAMODB_ENDPOINT=$(DYNAMODB_ENDPOINT) \
		TABLE_NAME_SUFFIX= \
		AWS_PROFILE=$(AWS_PROFILE) \
		./main || exit 0

generate:
	@go generate ./...

aws_account_id: ## current aws_account_id
	$(eval AWS_ACCOUNT_ID := $(shell $(AWS_VAULT_CMD) aws sts get-caller-identity --query 'Account' --output text))

help: ## Help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort |\
		awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
