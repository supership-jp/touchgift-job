PROJECT_NAME_PREFIX=touchgift-
DELIVERY_OPERATION_QUEUE = $(PROJECT_NAME_PREFIX)delivery-operation
DELIVERY_CONTROL_QUEUE = $(PROJECT_NAME_PREFIX)delivery-control

DELIVERY_OPERATION_TOPIC = $(PROJECT_NAME_PREFIX)delivery-operation-local
DELIVERY_OPERATION_TOPIC_ARN = arn:aws:sns:ap-northeast-1:000000000000:$(DELIVERY_OPERATION_TOPIC)
DELIVERY_OPERATION_NOTIFICATION_ENDPOINT = 'arn:aws:sqs:ap-northeast-1:000000000000:$(DELIVERY_OPERATION_QUEUE)'

DELIVERY_CONTROL_TOPIC = $(PROJECT_NAME_PREFIX)server-campaign-cache-local
DELIVERY_CONTROL_TOPIC_ARN = arn:aws:sns:ap-northeast-1:000000000000:$(DELIVERY_CONTROL_TOPIC)

start-localstack: ## build sqs mock (elasticmq)
	docker run --name touchgift-job-localstack -d --rm -it -p 4566:4566 -p 4567-4582:4567-4582 -e "LOCALSTACK_SERVICES=sns,sqs,dynamodb" localstack/localstack:3.7.1

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
create-sqs-delivery-operation-queue: ## delivery-operationキューを作成
	aws sqs create-queue --queue-name $(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)

create-sqs-delivery-control-queue: ## delivery-controlキューを作成
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

# キューへメッセージを送信する
send-message-sqs-delivery-operation-queue-insert-beta: ## delivery-operationキューへInsertメッセージ作成
	aws-vault exec tpo-staging -- aws sqs send-message --queue-url https://sqs.ap-northeast-1.amazonaws.com/411776208124/touchgift-delivery-operation-beta1 \
	--message-body $$(make show-queue-message LOG='{"time":"2024-09-01T09:52:57.804Z","type":"delivery_operation","request_id":"r1","target":"touchgift","campaigns":[{"id":1,"event":"insert","org_code":"org_1","creatives":[{"id":1,"event":"insert","org_code":"org_1"}]}]}')

send-message-sqs-delivery-operation-queue-insert:
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) \
	--message-body $$(make show-queue-message LOG='{"time":"2024-09-01T09:52:57.804Z","type":"delivery_operation","request_id":"r1","target":"touchgift","campaigns":[{"id":1,"event":"insert","org_code":"org_1","creatives":[{"id":1,"event":"insert","org_code":"org_1"}]}]}') $(SQS_OPTIONS)

send-message-sqs-delivery-operation-queue-update: ## delivery-operationキューへUpdateメッセージ作成
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) \
	--message-body $$(make show-queue-message LOG='{"time":"2024-09-01T09:52:57.804Z","type":"delivery_operation","request_id":"r2","target":"touchgift","campaigns":[{"id":1,"event":"update","org_code":"org_1","creatives":[{"id":1,"event":"update","org_code":"org_1"}]}]}') $(SQS_OPTIONS)

send-message-sqs-delivery-operation-queue-delete: ## delivery-operationキューへDeleteメッセージ作成
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) \
	--message-body $$(make show-queue-message LOG='{"time":"2024-09-01T09:52:57.804Z","type":"delivery_operation","request_id":"r3","target":"touchgift","campaigns":[{"id":1,"event":"delete","org_code":"org_1","creatives":[{"id":1,"event":"delete","org_code":"org_1"}]}]}') $(SQS_OPTIONS)

send-message-delivery-control: ##　delivery-controlキューへメッセージ送信
	aws sqs send-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) \
	--message-body $$(make show-queue-message LOG='{"time":"2021-08-02T09:42:43.388523304Z","version":1,"cache_operation":"DELETE","event":"expended","source":"touchgift-delivery-data-manager","org_code":"org_1","campaign_id":2}') $(SQS_OPTIONS)

receive-message-sqs-delivery-operation-queue: ## delivery-operationのメッセージ確認
	aws-vault exec tpo-staging -- aws sqs receive-message --queue-url https://sqs.ap-northeast-1.amazonaws.com/411776208124/touchgift-delivery-operation-beta1
receive-message-sqs-delivery-control-queue: ## delivery-controlのメッセージ確認
	aws sqs receive-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)


# キュー自体の削除
delete-sqs-delivery-operation-queue: ## delivery-operation-queueの削除
	@aws sqs delete-queue --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)
delete-sqs-delivery-control-queue: ## delivery-control-queueの削除
	@aws sqs delete-queue --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)

#　メッセージの削除
delete-message-sqs-delivery-operation-queue-beta: ## delivery-operationのメッセージをすべて削除
	aws-vault exec tpo-staging -- aws sqs purge-queue --queue-url https://sqs.ap-northeast-1.amazonaws.com/411776208124/touchgift-delivery-operation-beta1
delete-message-sqs-delivery-operation-queue: ## delivery-operationのメッセージをすべて削除
	aws sqs purge-queue --queue-url http://localhost:4566/000000000000/$(DELIVERY_OPERATION_QUEUE) $(SQS_OPTIONS)
delete-message-sqs-delivery-control-queue: ## delivery-controlのメッセージをすべて削除
	aws sqs receive-message --queue-url http://localhost:4566/000000000000/$(DELIVERY_CONTROL_QUEUE) $(SQS_OPTIONS)

