rm ./bin/main
cd service
go build -o ../bin/main
cd ../
export KITSVC_NAME="StringService"
export KITSVC_URL="http://127.0.0.1:8080"
export KITSVC_ADDR="127.0.0.1:8080"
export KITSVC_PORT=8080
export KITSVC_USAGE="Operations about the string."
export KITSVC_JWT_SECRET="4Rtg8BPKwixXy2ktDPxoMMAhRzmo9mmuZjvKONGPZZQSaJWNLijxR42qRgq0iBb5"
export KITSVC_VERSION="0.0.1"

export KITSVC_DATABASE_NAME="service"
export KITSVC_DATABASE_HOST="127.0.0.1:3306"
export KITSVC_DATABASE_USER="root"
export KITSVC_DATABASE_PASSWORD="root"
export KITSVC_DATABASE_CHARSET="utf8"
export KITSVC_DATABASE_LOC="Local"
export KITSVC_DATABASE_PARSE_TIME="True"

export KITSVC_ES_SERVER_URL="http://127.0.0.1:2113"
export KITSVC_ES_USERNAME="admin"
export KITSVC_ES_PASSWORD="changeit"

export KITSVC_PROMETHEUS_NAMESPACE="my_group"
export KITSVC_PROMETHEUS_NAMESPACE="string_service"
export KITSVC_CONSUL_CHECK_INTERVAL="30s"
export KITSVC_CONSUL_CHECK_TIMEOUT="1s"
export KITSVC_CONSUL_TAGS="string,micro"
./bin/main

#protoc --go_out=. ./pb/*.proto