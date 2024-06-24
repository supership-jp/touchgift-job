ifneq (,$(findstring beta,$(ENVIRONMENT_SUFFIX)))
	SNS_DELIVERY_OPERATION_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-operation-$(ENVIRONMENT_SUFFIX)
else
	SNS_DELIVERY_OPERATION_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-operation
endif
ifneq (,$(findstring rc,$(ENVIRONMENT_SUFFIX)))
	SNS_DELIVERY_OPERATION_LOG_TOPIC_ARN=arn:aws:sns:ap-northeast-1:$(AWS_ACCOUNT_ID):touchgift-operation-$(ENVIRONMENT_SUFFIX)
endif

create-delivery-operation-json: ## Create json for delivery operation. ex. make create-delivery-operation-json CAMPAIGN_ID=1 ORG=ss
	cat sns/sample.json \
		| jq ".campaigns[0].id|=$(CAMPAIGN_ID)" \
		| jq ".campaigns[0].organization_code|=\"$(ORG)\""

publish-delivery-operation: aws_account_id ## Publish delivery operation sns. ex. make publish-delivery-operation ENV=staging ENVIRONMENT_SUFFIX=beta1 PUBLISH_FILE=file://publish-temp.json
	$(AWS_VAULT_CMD) aws sns publish \
		--topic-arn $(SNS_DELIVERY_OPERATION_LOG_TOPIC_ARN) \
		--message $(PUBLISH_FILE)