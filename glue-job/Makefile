.DEFAULT_GOAL := help

CREDENTIAL = tpo-$(ENV)
AWS_VAULT_CMD = aws-vault exec $(CREDENTIAL) --
SUFFIX_WORKSPACE = '-$(WORKSPACE)'

JOB_NAME =
RUN_ID =

# Glueジョブをテストモードで実行
test-job: ##  make test-job WORKSPACE=beta JOB_NAME=touchgift-event-transfer ENV=staging
	$(AWS_VAULT_CMD) aws glue start-job-run --job-name '$(JOB_NAME)$(SUFFIX_WORKSPACE)' --arguments '{"--mode":"test"}'

# Glueジョブのテスト結果を表示
show-test-result: ##  make show-test-result WORKSPACE=beta JOB_NAME=touchgift-event-transfer RUN_ID=jr_e9cd37c4b9481ad6a6d17d1e5032d05a7c30975307e9af83a185ef1bad735a5d ENV=staging
	@id=$$($(AWS_VAULT_CMD) aws glue get-job-run --job-name '$(JOB_NAME)$(SUFFIX_WORKSPACE)' --run-id $(RUN_ID) | jq -r '.JobRun.Id'); \
	log_group_name=$$($(AWS_VAULT_CMD) aws glue get-job-run --job-name '$(JOB_NAME)$(SUFFIX_WORKSPACE)' --run-id $(RUN_ID) | jq -r '.JobRun.LogGroupName'); \
	full_log_group_name=$$log_group_name/output; \
	$(AWS_VAULT_CMD) aws logs get-log-events --log-group-name $$full_log_group_name --log-stream-name $$id

help: ## Help
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort |\
		awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'