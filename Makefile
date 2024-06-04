build:
	go build -o bin/eskimoe-server

doc:
	swag init --parseDependency --parseInternal

run:
	CompileDaemon -build="make build" -command="./bin/eskimoe-server"

test:
	go test -v ./...

clean:
	rm -rf bin