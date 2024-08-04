# 対象テーブル
## touch_pointテーブル
## campaignテーブル
## contentsテーブル
TN_PREFIX=
TABLE_NAME_SUFFIX=

TBL_TOUCH_POINT = $(TN_PREFIX)touchgift_delivery_data$(TABLE_NAME_SUFFIX)
TBL_CAMPAIGN = $(TN_PREFIX)touchgift_campaign_data$(TABLE_NAME_SUFFIX)
TBL_CONTENT = $(TN_PREFIX)touchgift_content_data$(TABLE_NAME_SUFFIX)
TBL_CREATIVE = $(TN_PREFIX)touchgift_creative_data$(TABLE_NAME_SUFFIX)

create-all-table: ## create all table
	$(MAKE) create-touch-point-table && \
	$(MAKE) create-campaign-table && \
	$(MAKE) create-content-table && \
	$(MAKE) create-creative-table
get-all-table: ## get all table
	$(MAKE) get-touch-point && \
	$(MAKE) get-campaign && \
	$(MAKE) get-contents && \
	$(MAKE) get-creatives
delete-all-table: ## delete all table
	$(MAKE) delete-touch-point && \
	$(MAKE) delete-campaign && \
	$(MAKE) delete-contents && \
	$(MAKE) delete-creatives

list-tables: ## dynamoのテーブルリスト一覧を表示します
	aws dynamodb list-tables $(DYNAMODB_OPTIONS)

create-touch-point-table: ## タッチポイントテーブルを作成する
	aws dynamodb create-table --table-name $(TBL_TOUCH_POINT) \
		--attribute-definitions \
			AttributeName=group_id,AttributeType=S \
			AttributeName=touchpoint_id,AttributeType=S \
		--key-schema \
			AttributeName=group_id,KeyType=HASH \
			AttributeName=touchpoint_id,KeyType=RANGE \
		--billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-touch_point: ## タッチポイント情報の取得
	aws dynamodb scan --table-name $(TBL_TOUCH_POINT) $(DYNAMODB_OPTIONS)
delete-touch_point: ## タッチポイントテーブルの削除
	aws dynamodb delete-table --table-name $(TBL_TOUCH_POINT) $(DYNAMODB_OPTIONS)

create-campaign-table: ## キャンペーンテーブルを作成する
	aws dynamodb create-table --table-name $(TBL_CAMPAIGN) \
		--attribute-definitions AttributeName=id,AttributeType=S \
		--key-schema AttributeName=id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-campaign: ## キャンペーン情報の取得
	aws dynamodb scan --table-name $(TBL_CAMPAIGN) $(DYNAMODB_OPTIONS)
delete-campaign: ## キャンペーンテーブルの削除
	aws dynamodb delete-table --table-name $(TBL_CAMPAIGN) $(DYNAMODB_OPTIONS)

create-content-table: ## コンテンツテーブルを作成する
	aws dynamodb create-table --table-name $(TBL_CONTENT) \
		--attribute-definitions AttributeName=campaign_id,AttributeType=S \
		--key-schema AttributeName=campaign_id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-content: ## コンテンツ情報の取得
	aws dynamodb scan --table-name $(TBL_CONTENTS) $(DYNAMODB_OPTIONS)
delete-content: ## コンテンツテーブルの削除
	aws dynamodb delete-table --table-name $(TBL_CONTENTS) $(DYNAMODB_OPTIONS)

create-creative-table: ## クリエイティブテーブルを作成する
	aws dynamodb create-table --table-name $(TBL_CREATIVE) \
		--attribute-definitions AttributeName=campaign_id,AttributeType=S \
		--key-schema AttributeName=campaign_id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-creatives:
	aws dynamodb scan --table-name $(TBL_CREATIVE) $(DYNAMODB_OPTIONS)
delete-creatives:
	aws dynamodb delete-table --table-name $(TBL_CREATIVE) $(DYNAMODB_OPTIONS)