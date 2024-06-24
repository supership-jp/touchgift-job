PROJECT_NAME_PREFIX=touchgift-
DELIVERY_OPERATION_QUEUE = $(PROJECT_NAME_PREFIX)delivery-operation
DELIVERY_CONTROL_QUEUE = $(PROJECT_NAME_PREFIX)delivery-control

DELIVERY_OPERATION_TOPIC = $(PROJECT_NAME_PREFIX)delivery-operation-local
DELIVERY_OPERATION_TOPIC_ARN = arn:aws:sns:ap-northeast-1:000000000000:$(DELIVERY_OPERATION_TOPIC)
DELIVERY_OPERATION_NOTIFICATION_ENDPOINT = '$(SNS_ENDPOINT)/000000000000/$(DELIVERY_OPERATION_QUEUE)'

DELIVERY_CONTROL_TOPIC = $(PROJECT_NAME_PREFIX)delivery-control-local
DELIVERY_CONTROL_TOPIC_ARN = arn:aws:sns:ap-northeast-1:000000000000:$(DELIVERY_CONTROL_TOPIC)

start-localstack: ## build sqs mock (elasticmq)
	docker run --name touchgift-job-localstack -d --rm -it -p 4566:4566 -p 4567-4582:4567-4582 -e "LOCALSTACK_SERVICES=sns,sqs,dynamodb" localstack/localstack:1.2.0

stop-localstack: ## stop sqs mock (elasticmq)
	docker stop touchgift-job-localstack
###
# SNS
###
create-sns-delivery-operation-topic: ## create delivery-operation topic at sns mock
	aws sns create-topic --name $(DELIVERY_OPERATION_TOPIC) $(SNS_OPTIONS)

create-sns-delivery-operation-subscribe: ## create delivery-operation subscribe at sns mock
	aws sns subscribe --topic-arn $(DELIVERY_OPERATION_TOPIC_ARN) --protocol sqs \
	--notification-endpoint $(DELIVERY_OPERATION_NOTIFICATION_ENDPOINT) $(SNS_OPTIONS)

delete-sns-delivery-operation-topic: ## delete delivery-operation topic at sns mock
	aws sns delete-topic --topic-arn $(DELIVERY_OPERATION_TOPIC_ARN) $(SNS_OPTIONS)

create-sns-delivery-control-topic: ## create delivery-control topic at sns mock
	aws sns create-topic --name $(DELIVERY_CONTROL_TOPIC) $(SNS_OPTIONS)

delete-sns-delivery-control-topic: ## delete delivery-control topic at sns mock
	aws sns delete-topic --topic-arn $(DELIVERY_CONTROL_TOPIC_ARN) $(SNS_OPTIONS)

create-sns-all: ## create sns and subscribe at sns mock
	$(MAKE) create-sns-delivery-operation-topic && \
	$(MAKE) create-sns-delivery-control-topic && \
	$(MAKE) create-sns-delivery-operation-subscribe

list-sns-topics: ## list topics at sns mock
	@aws sns list-subscriptions-by-topic --topic-arn $(DELIVERY_OPERATION_TOPIC_ARN) $(SNS_OPTIONS)

list-sns-subscriptions: ## list topics at sns mock
	@aws sns list-subscriptions $(SNS_OPTIONS)

# SUBSCRIPTION_ARN: list-sns-subscriptionsから取得したSubscriptionArn
sns-unsubscribe: ## sns unsubscribe ex: make sns-unsubscribe SUBSCRIPTION_ARN=
	@aws sns unsubscribe --subscription-arn $(SUBSCRIPTION_ARN) $(SNS_OPTIONS)

# ここからローカル用
list-all-topics:
	 @aws sns list-topics --endpoint-url=http://localhost:4566 --profile=dummy
list-all-subscriptions:
	@aws sns list-subscriptions --endpoint-url=http://localhost:4566 --profile=dummy