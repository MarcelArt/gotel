swag:
	@swag init --parseDependency --parseInternal

dev: swag
	@air serve

go: swag
	@go run main.go serve

migrate:
	@go run main.go migrate

migrate-force:
	@go run main.go migrate --drop

docker-build:
	@docker build -t marcelaritonang/gotel:latest .

docker-run: docker-build
	@docker run -p 7001:7001 --env-file .env marcelaritonang/gotel:latest

docker-push: docker-build
	@docker push marcelaritonang/gotel:latest