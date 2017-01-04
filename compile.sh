rm ./bin/KitSvc
cd service
go build -o ../bin/KitSvc
cd ../
export KITSVC_NAME="StringService"
export KITSVC_URL="http://127.0.0.1:8080"
export KITSVC_ADDR="127.0.0.1:8080"
export KITSVC_PORT=8080
export KITSVC_USAGE="Operations about the string."
export KITSVC_VERSION="0.0.1"

export KITSVC_DATABASE_NAME="service"
export KITSVC_DATABASE_HOST="127.0.0.1:3306"
export KITSVC_DATABASE_USER="root"
export KITSVC_DATABASE_PASSWORD="root"
export KITSVC_DATABASE_CHARSET="utf8"
export KITSVC_DATABASE_LOC="Local"
export KITSVC_DATABASE_PARSE_TIME="True"

export KITSVC_NSQ_PRODUCER="127.0.0.1:4150"
export KITSVC_NSQ_LOOKUPS="127.0.0.1:4161"

export KITSVC_PROMETHEUS_NAMESPACE="my_group"
export KITSVC_PROMETHEUS_NAMESPACE="string_service"
export KITSVC_CONSUL_CHECK_INTERVAL="10s"
export KITSVC_CONSUL_CHECK_TIMEOUT="1s"
export KITSVC_CONSUL_TAGS="string,micro"
./bin/KitSvc