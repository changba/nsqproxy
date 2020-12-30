GO111MODULE=on

filename = nsqproxy
nowDate = $(shell date +"%Y%m%d%H%M%S")

.PHONY: build
build:
	@echo "Build..."
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/$(filename) cmd/nsqproxy.go

.PHONY: build-linux
build-linux:
	@echo "Build for linux..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(filename)-linux-$(nowDate) cmd/nsqproxy.go

.PHONY: build-all
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(filename)-linux-$(nowDate) cmd/nsqproxy.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/$(filename)-darwin-$(nowDate) cmd/nsqproxy.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(filename)-windows-$(nowDate) cmd/nsqproxy.go

.PHONY: clean
clean:
	rm bin/nsqproxy*

.PHONY: kill
kill:
	-killall nsqproxy

.PHONY: test
test:
	go test ./...

.PHONY: run
run: build kill
	./bin/nsqproxy &

.PHONY: statik
statik:
	go get github.com/rakyll/statik
	go generate ./...

.PHONY: vue-install
vue-install:
	cd web/vue-admin && npm install

.PHONY: vue-install-taobao
vue-install-taobao:
	cd web/vue-admin && npm install --registry=https://registry.npm.taobao.org

.PHONY: vue-build
vue-build:
	cd web/vue-admin && npm run build:prod
	mkdir -p web/public && cp -r web/vue-admin/dist/* web/public/

.PHONY: vue-dev
vue-dev:
	cd web/vue-admin && yarn run dev