cd ./service
go tool cover -html=cover.out
go test -v -cover -tags test -coverprofile=cover.out