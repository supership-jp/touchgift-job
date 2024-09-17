.DEFAULT_GOAL := help
SHELL=/bin/bash

APPLICATION = touchgift-job-fluentd-daemon
ECS_CLUSTER = touchgift-job
SERVICE = $(APPLICATION)

CURRENT_TAG = 1.0.0
TAG = $(CURRENT_TAG)
ifndef TAG
override TAG = $(CURRENT_TAG)
endif
REPOSITORY = touchgift/job-fluentd-daemon
FLUENTD_TAG = v1.13.3-debian-1.0

CREDENTIAL = tpo-$(ENV)
AWS_VAULT_CMD = aws-vault exec $(CREDENTIAL) --

AWS_REGION=ap-northeast-1

# e1, beta1 ...
ENVIRONMENT_SUFFIX = e1

build-image: aws_account_id ## Build docker image ex. make build-image ENV=staging TAG=
	@docker build --tag $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)-$(FLUENTD_TAG) .

push-image: aws_account_id ## Push docker image ex. make push-image ENV=staging TAG=
	@$(AWS_VAULT_CMD) aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com && \
	docker push $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)-$(FLUENTD_TAG) && \
	docker logout https://$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

output-image: aws_account_id ## output full image name make output-image CREDENTIAL=tpo-staging
	@echo $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)-$(FLUENTD_TAG)

run-deploy: ## run deploy ex. make run-deploy ENV=staging ENVIRONMENT_SUFFIX= TAG=
	@TAG=$(TAG) AWS_VAULT_CMD="$(AWS_VAULT_CMD)" ENV=$(ENV) ENVIRONMENT_SUFFIX=$(ENVIRONMENT_SUFFIX) ECS_CLUSTER=$(ECS_CLUSTER) SERVICE=$(SERVICE) ./run-deploy-task

register-task-definition: ## register task definition ex. make register-task-definition ENV=staging TASK_DEF_FILE=file://./task-definition.json
	@$(AWS_VAULT_CMD) aws ecs register-task-definition --cli-input-json $(TASK_DEF_FILE)

rewrite-image-on-task-definition: aws_account_id ## rewrite image on task-definition.json ex. make rewrite-image-on-task-definition TAG= ENV= ENVIRONMENT_SUFFIX=
	@cat task-definitions/$(ENV).json \
		| jq ".family|=\"$(SERVICE)-$(ENVIRONMENT_SUFFIX)\"" \
		| jq ".taskRoleArn|=\"arn:aws:iam::$(AWS_ACCOUNT_ID):role/$(APPLICATION)-$(ENVIRONMENT_SUFFIX).ecs-task-role\"" \
		| jq ".executionRoleArn|=\"arn:aws:iam::$(AWS_ACCOUNT_ID):role/$(APPLICATION)-$(ENVIRONMENT_SUFFIX).ecs-task-execution-role\"" \
		| jq ".containerDefinitions[0].image|=\"$(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(REPOSITORY):$(TAG)-$(FLUENTD_TAG)\"" \
		| jq ".containerDefinitions[0].logConfiguration.options[\"awslogs-group\"]|=\"/$(ECS_CLUSTER)-$(ENVIRONMENT_SUFFIX)/$(SERVICE)-$(ENVIRONMENT_SUFFIX)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"ENVIRONMENT_SUFFIX\") | .value) |=\"$(ENVIRONMENT_SUFFIX)\"" \
		| jq "(.containerDefinitions[0].environment[] | select(.name==\"FLUENTD_RECEIVER_SD\") | .value) |=\"fluentd-receiver-$(ENVIRONMENT_SUFFIX)\""

update-service: ## update service ex. make update-service TASK_DEF_ARN= ECS_CLUSTER=
	@$(AWS_VAULT_CMD) aws ecs update-service --cluster $(ECS_CLUSTER)-$(ENVIRONMENT_SUFFIX) --service $(SERVICE)-$(ENVIRONMENT_SUFFIX) --task-definition $(TASK_DEF_ARN)

help: ## Help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort |\
		awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

aws_account_id: ## current aws_account_id
	$(eval AWS_ACCOUNT_ID := $(shell $(AWS_VAULT_CMD) aws sts get-caller-identity --query 'Account' --output text))
