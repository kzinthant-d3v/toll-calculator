gate:
	@go build -o bin/gate gateway/main.go
	@./bin/gate
obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

producer:
	@go build -o bin/producer ./producer
	@./bin/producer

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

invoice:
	@go build -o bin/invoice ./invoice
	@./bin/invoice

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/*.proto

.PHONY: obu producer invoice
