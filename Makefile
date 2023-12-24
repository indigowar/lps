gen:
	templ generate ./...
	go generate ./...

run: gen
	@go run cmd/lps/main.go


up: gen
	docker-compose -f build/local-compose.yaml --env-file .env up --build --remove-orphans