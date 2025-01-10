build : clean
	@go build -o bin/api cmd/api/main.go
clean:
	@rm -rf bin
db_create_migration :
	migrate create -ext sql -dir migrations -seq $(name)
db_migrate:
	# Load the environment variables from the .env file and export them
	@set -a && source .env && set +a && \
		migrate -database $$DATABASE_URL -path migrations up
db_login:
	@set -a && source .env && set +a && \
		USER=$$DB_USER && \
		PASSWORD=$$DB_PASS && \
		HOST=$$DB_HOST && \
		PORT=$$DB_PORT && \
		mysql -u $$USER -p$$PASSWORD -h $$HOST -P $$PORT
