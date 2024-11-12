include ./localstack/Makefile
include ./dynamodb/Makefile
include ./sns/Makefile

.DEFAULT_GOAL := help
SHELL=/bin/bash

APPLICATION = touchgift-job
ECS_CLUSTER = touchgift-job
SERVICE = $(APPLICATION)

AWS_PROFILE = dummy
AWS_REGION=ap-northeast-1

# Dynamo設定まわり
DYNAMODB_ENDPOINT = http://localhost:8000
DYNAMODB_ENDPOINT_OPTS = --endpoint-url $(DYNAMODB_ENDPOINT)
DYNAMODB_OPTIONS = $(DYNAMODB_ENDPOINT_OPTS) --profile $(AWS_PROFILE)

CURRENT_TAG = 1.0.1
TAG = $(CURRENT_TAG)
ifndef TAG
override TAG = $(CURRENT_TAG)
endif

TEST_TARGET=./...

REPOSITORY = touchgift/job
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
### e1(本番), rc1(リリース前テスト環境)
## 検証アカウント
### e1(検証), beta1(ベータ環境)
ENVIRONMENT_SUFFIX = e1

# suffixで環境変数を切り替える
ifneq (,$(findstring beta,$(ENVIRONMENT_SUFFIX)))
	DELIVERY_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CAMPAIGN_TABLE_NAME=touchgift_campaign_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=touchgift_creative_data.$(ENVIRONMENT_SUFFIX)
	TOUCH_POINT_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CONTENT_TABLE_NAME=touchgift_content_data.$(ENVIRONMENT_SUFFIX)
	DB_HOST=/rds/touchgift-main-beta/endpoint
	DB_PASSWORD=/rds/touchgift-main-beta/user/sc2_pn/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-campaign-cache-$(ENVIRONMENT_SUFFIX)
	SNS_DELIVERY_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-delivery-cache-$(ENVIRONMENT_SUFFIX)
	SNS_CREATIVE_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-creative-cache-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/touchgift-delivery-operation-$(ENVIRONMENT_SUFFIX)
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/touchgift-delivery-control-$(ENVIRONMENT_SUFFIX)
else
	TOUCH_POINT_TABLE_NAME=touchgift_delivery_data
	CREATIVE_TABLE_NAME=touchgift_creative_data
	CAMPAIGN_TABLE_NAME=touchgift_campaign_data
	CONTENT_TABLE_NAME=touchgift_content_data
	DB_HOST=/rds/touchgift-main/endpoint
	DB_PASSWORD=/rds/touchgift-main/user/sc2_pn/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-campaign-cache
	SNS_DELIVERY_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-delivery-cache
	SNS_CREATIVE_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-creative-cache
	SQS_DELIVERY_OPERATION_QUEUE_URL=/sqs/url/touchgift-delivery-operation
	SQS_DELIVERY_CONTROL_QUEUE_URL=/sqs/url/touchgift-delivery-control
endif
ifneq (,$(findstring rc,$(ENVIRONMENT_SUFFIX)))
	DELIVERY_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CAMPAIGN_TABLE_NAME=touchgift_campaign_data.$(ENVIRONMENT_SUFFIX)
	CREATIVE_TABLE_NAME=touchgift_creative_data.$(ENVIRONMENT_SUFFIX)
	TOUCH_POINT_TABLE_NAME=touchgift_delivery_data.$(ENVIRONMENT_SUFFIX)
	CONTENT_TABLE_NAME=touchgift_content_data.$(ENVIRONMENT_SUFFIX)
	DB_HOST=/rds/touchgift-main-rc/endpoint
	DB_PASSWORD=/rds/touchgift-main-rc/user/sc2_pn/password
	SNS_CONTROL_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-campaign-cache-$(ENVIRONMENT_SUFFIX)
	SNS_DELIVERY_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-delivery-cache-$(ENVIRONMENT_SUFFIX)
	SNS_CREATIVE_CACHE_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-server-creative-cache-$(ENVIRONMENT_SUFFIX)
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
		-e MYSQL_DATABASE=retail \
		-e MYSQL_USER=sc2_pn \
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

run-deploy: ## run deploy (code deploy) ex. make run-deploy CREDENTIAL=tpo-staging TAG=testtest ENVIRONMENT= ENVIRONMENT_SUFFIX=
	@REGISTRY=$(REGISTRY) TAG=$(TAG) CREDENTIAL=$(CREDENTIAL) AWS_VAULT_CMD="$(AWS_VAULT_CMD)" \
		ENVIRONMENT=$(ENVIRONMENT) ECS_CLUSTER=$(ECS_CLUSTER) ENVIRONMENT_SUFFIX=$(ENVIRONMENT_SUFFIX) ./run-deploy-task

register-task-definition: ## register task definition ex. make register-task-definition ENV=staging TASK_DEF_FILE=file://./task-definition.json
	@$(AWS_VAULT_CMD) aws ecs register-task-definition --cli-input-json $(TASK_DEF_FILE)

rewrite-image-on-task-definition:aws_account_id ## rewrite image on task-definition.json ex. make rewrite-image-on-task-definition ENV= ENVIRONMENT_SUFFIX= TAG=
	@cat task-definitions/$(ENV).json \
		| jq ".family|=\"$(SERVICE)-$(ENVIRONMENT_SUFFIX)\"" \
		| jq ".taskRoleArn|=\"arn:aws:iam::$(AWS_ACCOUNT_ID):role/$(APPLICATION)-$(ENVIRONMENT_SUFFIX).ecs-task-role\"" \
		| jq ".executionRoleArn|=\"arn:aws:iam::$(AWS_ACCOUNT_ID):role/$(APPLICATION)-$(ENVIRONMENT_SUFFIX).ecs-task-execution-role\"" \
		| jq ".containerDefinitions[0].image|=\"$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)\"" \
		| jq ".containerDefinitions[0].logConfiguration.options.tag|=\"touchgiftjob.$(ENVIRONMENT_SUFFIX).$(subst .,-,$(TAG)).{{.ID}}\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"ENVIRONMENT_SUFFIX\") | .value) |=\"$(ENVIRONMENT_SUFFIX)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"AWS_REGION\") | .value) |=\"$(AWS_REGION)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"CAMPAIGN_TABLE_NAME\") | .value) |=\"$(CAMPAIGN_TABLE_NAME)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"CREATIVE_TABLE_NAME\") | .value) |=\"$(CREATIVE_TABLE_NAME)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"TOUCH_POINT_TABLE_NAME\") | .value) |=\"$(TOUCH_POINT_TABLE_NAME)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"CONTENT_TABLE_NAME\") | .value) |=\"$(CONTENT_TABLE_NAME)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"SNS_CONTROL_LOG_TOPIC_ARN\") | .value) |=\"$(SNS_CONTROL_LOG_TOPIC_ARN)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"SNS_DELIVERY_CACHE_TOPIC_ARN\") | .value) |=\"$(SNS_DELIVERY_CACHE_TOPIC_ARN)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"SNS_CREATIVE_CACHE_TOPIC_ARN\") | .value) |=\"$(SNS_CREATIVE_CACHE_TOPIC_ARN)\"" \
		| jq "(.containerDefinitions[0].secrets[] | select(.name==\"DB_HOST\") | .valueFrom) |=\"$(DB_HOST)\"" \
		| jq "(.containerDefinitions[0].secrets[] | select(.name==\"DB_PASSWORD\") | .valueFrom) |=\"$(DB_PASSWORD)\"" \
		| jq "(.containerDefinitions[0].secrets[] | select(.name==\"SQS_DELIVERY_OPERATION_QUEUE_URL\") | .valueFrom) |=\"$(SQS_DELIVERY_OPERATION_QUEUE_URL)\"" \
		| jq "(.containerDefinitions[0].secrets[] | select(.name==\"SQS_DELIVERY_CONTROL_QUEUE_URL\") | .valueFrom) |=\"$(SQS_DELIVERY_CONTROL_QUEUE_URL)\""

update-service: ## update service ex. make update-service ENV= ENVIRONMENT_SUFFIX= TASK_DEF_ARN=
	@$(AWS_VAULT_CMD) aws ecs update-service --cluster $(ECS_CLUSTER)-$(ENVIRONMENT_SUFFIX) --service $(SERVICE)-$(ENVIRONMENT_SUFFIX) --task-definition $(TASK_DEF_ARN)

update-desired-count: ## update ecs desiredCount ex. make update-desired-count ENV=staging ENVIRONMENT_SUFFIX=beta1 COUNT=1
	@$(AWS_VAULT_CMD) aws ecs update-service \
		--cluster $(ECS_CLUSTER)-$(ENVIRONMENT_SUFFIX) \
		--service $(SERVICE)-$(ENVIRONMENT_SUFFIX) \
		--desired-count $(COUNT)
