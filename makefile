clean: 
	go clean
build:
	go build -o bin/zomboidBot
docker-build:
	docker build -t zomboid-bot-image .
