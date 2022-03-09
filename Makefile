.PHONY: test
test: api-test ui-test ## run all Unit Tests

.PHONY: lint
lint: api-lint ui-lint yaml-lint ## check all Lints

.PHONY: api-check
api-check: goa-gen api-test generated api-lint ## all API checks

.PHONY: ui-check
ui-check: ui-test ui-lint ## all UI checks

.PHONY: goa-gen
goa-gen: ## generate API Design
	@echo "----------------------------"
	@echo "-- Generating API Design... "
	@echo "----------------------------"
	cd api && go mod vendor && goa gen github.com/tektoncd/hub/api/design

	@echo "----------------------------"
	@echo "- Generating v1 API Design... "
	@echo "----------------------------"
	cd api/v1 && go mod vendor && goa gen github.com/tektoncd/hub/api/v1/design

.PHONY: generated
generated: ## update the golden files
	@echo "----------------------------"
	@echo "-- Generating Golden Files... --"
	@echo "----------------------------"
	cd api && go mod vendor && bash update-golden-files.sh

.PHONY: api-lint
api-lint: ## run API Lint
	@echo "----------------------------"
	@echo "-- Linting API Files... --"
	@echo "----------------------------"
	cd api && golangci-lint run -v ./pkg/... ./v1/service/... --timeout=5m


.PHONY: api-test
api-test: ## run API Unit Test
	@echo "----------------------------"
	@echo "-- Running API Unit Tests --"
	@echo "----------------------------"
	@echo "🛢 🛢 🛢  Make sure you have your Database Up and Running 🛢 🛢 🛢"
	@echo "----------------------------"
	cd api && go mod vendor && go test -p 1 -v ./pkg/... ./v1/service/...

.PHONY: api-build
api-build: ## generate the API binary
	@echo "----------------------------"
	@echo "-- Running API Build... --"
	@echo "----------------------------"
	cd api && go mod vendor && go build -mod=vendor ./cmd/...


.PHONY: ui-lint
ui-lint: ## run UI Lint
	@echo "----------------------------"
	@echo "-- Linting UI Files... --"
	@echo "----------------------------"
	cd ui && npm clean-install && npm run lint

.PHONY: ui-test
ui-test: ## run UI Unit Test
	@echo "----------------------------"
	@echo "-- Running UI Unit Tests --"
	@echo "----------------------------"
	cd ui && npm clean-install && CI=true npm test

.PHONY: ui-build
ui-build: ## generate the UI binary
	@echo "----------------------------"
	@echo "-- Running UI Build... --"
	@echo "----------------------------"
	cd ui && npm clean-install && CI=true npm run build

.PHONY: yaml-lint
yaml-lint: ## run YAML Lint
	@echo "----------------------------"
	@echo "-- Running Yaml-lint... --"
	@echo "----------------------------"
	yamllint -c .yamllint ./config.yaml ./config

.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
