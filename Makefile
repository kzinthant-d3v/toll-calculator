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

.PHONY: obu, producer, invoice
