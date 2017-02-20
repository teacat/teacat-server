.PHONY: build

deps:
	cp -a vendor/. $(GOPATH)/src

test:
	go test ./server

test_sqlite:
	go test github.com/TeaMeow/KitSvc/store/datastore

test_mysql:
	KITSVC_DATABASE_DRIVER="mysql" KITSVC_DATABASE_CONFIG="root:root@tcp(localhost:3306)/service?charset=utf8&parseTime=True&loc=Local" go test github.com/TeaMeow/KitSvc/store/datastore

run:
	go build -o ./bin/main ./server
	./bin/main

build:
	go build -o ./bin/main ./server