DELIVERY_OPERATION_QUEUE = delivery-operation
DELIVERY_CONTROL_QUEUE = delivery-control

start-localstack: ## build sqs mock (elasticmq)
	docker run --name touchgift-localstack -d --rm -it -p 4566:4566 -p 4567-4582:4567-4582 -e "LOCALSTACK_SERVICES=sns,sqs,dynamodb" localstack/localstack:1.2.0

stop-localstack: ## stop sqs mock (elasticmq)
	docker stop adm-ddm-localstack