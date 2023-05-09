build:
	go build -o bin/distributed_cache

run: build
	./bin/distributed_cache