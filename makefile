clean: 
	go clean
build:
	go build
docker-build:
	docker build -t zomboid-bot-image .
