KV_STORE_RELEASE ?= dev
build:
	go build -o ./.bin/kv-store -ldflags "-X github.com/exsql-io/kv-store/pkg/kvstore.Version=$(KV_STORE_RELEASE)"  cmd/kvstore/kv_store.go