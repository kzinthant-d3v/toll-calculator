obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./receiver
	@./bin/receiver

.PHONY: obu
.PHONY: receiver
