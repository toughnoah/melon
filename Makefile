TAG = v1.0
BUILD=webhook


run:
	@go run cmd/main.go

build:
	@export CGO_ENABLED=0 && go build -o $(BUILD) cmd/main.go

image: build
	@echo $(TAG)
	@docker build . -t toughnoah/melon:$(TAG)
	@docker push toughnoah/melon:$(TAG)
	@rm -f ./$(BUILD)

test:
	@go test ./... -v -coverprofile=cover.out
	@go tool cover -func=cover.out

clean:
	#clean useless image
	@docker rm `docker ps -a | grep Exited | awk '{print $1}'`
	@docker rmi -f  `docker images | grep '<none>' | awk '{print $3}'`