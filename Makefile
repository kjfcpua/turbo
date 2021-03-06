.PHONY: install
install:
	@cd turbo && go install
	@cd protoc-gen-buildfields && go install

.PHONY: test
test:
	@go test -cover -coverpkg github.com/vaporz/turbo github.com/vaporz/turbo github.com/vaporz/turbo/test

.PHONY: doc
doc:
	@cd doc && make html
