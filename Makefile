gqlgen:
	cd internal/delivery/graphqlsvc/schema && \
	go get github.com/99designs/gqlgen@v0.17.22 && \
	go run github.com/99designs/gqlgen generate

check-cognitive-complexity:
	find . -type f -name '*.go' -not -name "*.pb.go" -not -name "mock*.go" -not -name "generated.go" -not -name "federation.go" \
      -exec gocognit -over 15 {} +

lint: check-cognitive-complexity
	golangci-lint run --print-issued-lines=false --exclude-use-default=false --enable=revive --enable=goimports  --enable=unconvert --enable=unparam --concurrency=2 --skip-dirs='generated'

run: check-modd-exists
	@modd -f ./.modd/server.modd.conf

check-modd-exists:
	@modd --version > /dev/null

check-gotest:
ifeq (, $(shell which richgo))
	$(warning "richgo is not installed, falling back to plain go test")
	$(eval TEST_BIN=go test)
else
	$(eval TEST_BIN=richgo test)
endif

ifdef test_run
	$(eval TEST_ARGS := -run $(test_run))
endif
	$(eval test_command=$(TEST_BIN) ./... $(TEST_ARGS) -v --cover)

internal/model/mock/mock_composer_usecase.go:
	mockgen -destination=internal/model/mock/mock_composer_usecase.go -package=mock github.com/atjhoendz/notpushcation-service/internal/model ComposerUsecase

internal/model/mock/mock_message_processor_usecase.go:
	mockgen -destination=internal/model/mock/mock_message_processor_usecase.go -package=mock github.com/atjhoendz/notpushcation-service/internal/model MessageProcessorUsecase

internal/model/mock/mock_onesignal_client.go:
	mockgen -destination=internal/model/mock/mock_onesignal_client.go -package=mock github.com/atjhoendz/notpushcation-service/internal/model OnesignalClient

mockgen: internal/model/mock/mock_composer_usecase.go \
	internal/model/mock/mock_message_processor_usecase.go \
	internal/model/mock/mock_onesignal_client.go

test-only: check-gotest mockgen
	SVC_DISABLE_CACHING=true $(test_command) -timeout 60s

test: lint test-only

clean:
	rm -v internal/model/mock/mock_*.go

.PHONY: gqlgen check-cognitive-complexity lint