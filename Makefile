## Format source code
format:
	go fmt ./...

generate_docs:
	swag init

stop:
	docker-compose down -v

logs:
	docker-compose logs -f -t

clean:
	chmod -R +w ./.gopath vendor

##Runs application in docker container
start:
	docker-compose up -d

