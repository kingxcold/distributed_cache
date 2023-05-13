build:
	go build -o bin/distributed_cache

run: build
	./bin/distributed_cache

runfollower: build
	./bin/distributed_cache --listenaddr :4000 --leaderaddr :3000

test:
	go test