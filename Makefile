generate:
	@protoc --proto_path=proto --go_out=proto/gen --go-grpc_out=proto/gen --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/user.proto

build:
	@echo "---- Building Application ----"
	@go build -o server main.go

air_init:
	@air init

air:
	@air