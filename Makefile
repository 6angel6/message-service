
dc:
	docker-compose up --remove-orphans --build
run:
	go build -o app cmd/main.go
test:
	go test -race ./..