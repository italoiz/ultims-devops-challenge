BINARY   := api
IMAGE    := italoiz/ultims-event-tracker
TAG      := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test run lint docker-build docker-push clean

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BINARY) ./cmd/api

test:
	go test -v -race -count=1 ./...

run: build
	./$(BINARY)

lint:
	go vet ./...

docker-build:
	docker build -t $(IMAGE):$(TAG) -t $(IMAGE):latest .

docker-push: docker-build
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

clean:
	rm -f $(BINARY) coverage.out
