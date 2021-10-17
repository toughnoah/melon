TAG = v1.0

run:
	@go run cmd/main.go

build:
	@export CGO_ENABLED=0 &&
	@go build -o webhook cmd/main.go

image: build
	@echo $(TAG)
	@docker build . -t toughnoah/melon:$(TAG)
	@docker push toughnoah/melon:$(TAG)