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

###
# SQS
###
create-sqs-delivery-operation-queue: ## create queue at sqs mock (elasticmq)
	aws sqs create-queue --queue-name $(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)

create-sqs-delivery-control-queue: ## create queue at sqs mock (elasticmq)
	aws sqs create-queue --queue-name $(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)

create-sqs-all-queue: ## create queue at sqs mock (elasticmq)
	$(MAKE) create-sqs-delivery-operation-queue && \
	$(MAKE) create-sqs-delivery-control-queue

list-sqs-queue: ## list queue at sqs mock (elasticmq)
	aws sqs list-queues $(SQS_OPTIONS)

show-queue-message: ## show queue message
	@echo "{\"Type\":\"Notification\",\"MessageId\":\"id\",\"TopicArn\":\"arn\",\"Subject\":\"Notification\",\"Message\":\"$${LOG//\"/\\\"}\",\"Timestamp\":\"2021-02-16T07:35:25.158Z\",\"SignatureVersion\":\"1\",\"Signature\":\"Signature\",\"SigningCertURL\":\"\",\"UnsubscribeURL\":\"\"}"

test-a-a:
	@echo "$$(make show-queue-message LOG='{"a":"b"}')"


send-message-sqs-delivery-operation-queue-insert: ## delivery-operationキューへInsertメッセージ作成
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) --message-body $$(make show-queue-message LOG='{"time":"2020-03-04T09:52:57.804Z","type":"delivery_operation","request_id":"r1","target":"touchgift","schedules":[{"id":1,"origin_id":"origin_id1","budget":0,"event":"insert","organization_code":"code1","service":"default","advertiser_id":1,"locations":[],"creatives":[]}]}') $(SQS_OPTIONS)

send-message-sqs-delivery-operation-queue-update: ## delivery-operationキューへUpdateメッセージ作成
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) --message-body $$(make show-queue-message LOG='{"time":"2020-03-04T09:52:57.804Z","type":"delivery_operation","request_id":"r2","target":"touchgift","schedules":[{"id":2,"origin_id":"origin_id1","budget":0,"event":"update","organization_code":"code1","service":"default","advertiser_id":1,"locations":[],"creatives":[]}]}') $(SQS_OPTIONS)

send-message-sqs-delivery-operation-queue-delete: ## delivery-operationキューへDeleteメッセージ作成
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) --message-body $$(make show-queue-message LOG='{"time":"2020-03-04T09:52:57.804Z","type":"delivery_operation","request_id":"r3","target":"touchgift","schedules":[{"id":1,"origin_id":"origin_id1","budget":0,"event":"delete","organization_code":"code1","service":"default","advertiser_id":1,"locations":[],"creatives":[{"id":1,"event":"delete","organization_code":"o1","service":"default","advertiser_id":1}]}]}') $(SQS_OPTIONS)

send-message-delivery-control: ##　delivery-operationキューへメッセージ送信
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) --message-body $$(make show-queue-message LOG='{"time":"2021-08-02T09:42:43.388523304Z","version":1,"cache_operation":"DELETE","event":"expended","source":"touchgift-delivery-data-manager","organization":"code1","service":"default","schedule_id":2}') $(SQS_OPTIONS)

receive-message-sqs-delivery-operation-queue: ## delivery-operationのメッセージ確認
	aws sqs receive-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)
receive-message-sqs-delivery-control-queue:　## delivery-controlのメッセージ確認
	aws sqs receive-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)

delete-message-sqs-delivery-operation-queue: ## delivery-operationのメッセージをすべて削除
	aws sqs purge-queue --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)
delete-message-sqs-delivery-control-queue:　## delivery-controlのメッセージをすべて削除
	aws sqs receive-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)
