# serverless create -t aws-go-dep -p github.com/skyzyx/lambda-htmlmeta
all:
	@cat Makefile | grep : | grep -v PHONY | grep -v @ | sed 's/:/ /' | awk '{print $$1}' | sort

#-------------------------------------------------------------------------------

.PHONY: install-deps
install-deps:
	gometalinter.v2 --install

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/htmlinfo main.go

.PHONY: package
package:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/htmlinfo main.go

.PHONY: lint
lint:
	gometalinter.v2 ./main.go

#-------------------------------------------------------------------------------

.PHONY: deploy
deploy: package
	sls deploy --verbose
