.PHONY: help
.DEFAULT_GOAL := help
environment = "example"

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

create: ## create env
	@sceptre launch-env $(environment)

delete: ## delete env
	@sceptre delete-env $(environment)

info: ## describe resources
	@sceptre describe-stack-outputs $(environment) datasegment

publish: ## publish records
	$(shell KINESIS_STREAM_NAME=`sceptre --output json describe-stack-outputs example datasegment | jq -r '.[] | select(.OutputKey=="KinesisStreamName") | .OutputValue'` go run publisher.go)

subscribe: ## subscribe
	$(shell KINESIS_STREAM_NAME=`sceptre --output json describe-stack-outputs example datasegment | jq -r '.[] | select(.OutputKey=="KinesisStreamName") | .OutputValue'` go run subscriber.go)