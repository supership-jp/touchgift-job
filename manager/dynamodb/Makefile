# 対象テーブル
## touch_pointテーブル
## campaignテーブル
## contentsテーブル
TN_PREFIX=

TBL_TOUCH_POINT = $(TN_PREFIX)touch_point
TBL_CAMPAIGN = $(TN_PREFIX)campaign
TBL_CONTENTS= $(TN_PREFIX)contents

create-all-table: ## create all table
	$(MAKE) create-touch_point-table && \
	$(MAKE) create-campaign-table && \
	$(MAKE) create-contents-table
get-all-table: ## get all table
	$(MAKE) get-touch_point && \
	$(MAKE) get-campaign && \
	$(MAKE) get-content
delete-all-table: ## delete all table
	$(MAKE) delete-touch_point && \
	$(MAKE) delete-campaign && \
	$(MAKE) delete-contents

list-tables: ## dynamoのテーブルリスト一覧を表示します
	aws dynamodb list-tables $(DYNAMODB_OPTIONS)

create-touch_point-table: ## タッチポイントテーブルを作成する
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
		--attribute-definitions AttributeName=id,AttributeType=N \
		--key-schema AttributeName=id,KeyType=HASH \
		--billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-campaign: ## キャンペーン情報の取得
	aws dynamodb scan --table-name $(TBL_CAMPAIGN) $(DYNAMODB_OPTIONS)
delete-campaign: ## キャンペーンテーブルの削除
	aws dynamodb delete-table --table-name $(TBL_CAMPAIGN) $(DYNAMODB_OPTIONS)

create-contents-table: ## コンテンツテーブルを作成する
	aws dynamodb create-table --table-name $(TBL_CONTENTS) \
        --attribute-definitions AttributeName=campaign_id,AttributeType=S \
        --key-schema AttributeName=campaign_id,KeyType=HASH \
        --billing-mode PAY_PER_REQUEST $(DYNAMODB_OPTIONS)
get-contents: ## コンテンツ情報の取得
	aws dynamodb scan --table-name $(TBL_CONTENTS) $(DYNAMODB_OPTIONS)
delete-contents: ## コンテンツテーブルの削除
	aws dynamodb delete-table --table-name $(TBL_CONTENTS) $(DYNAMODB_OPTIONS)