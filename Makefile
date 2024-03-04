migrate-create:
	goose -dir ./migrates postgres "host=localhost port=5433 user=postgres dbname=rental sslmode=disable" create $(NAME) sql

migrate-up:
	@goose -dir ./migrates postgres "host=localhost port=5433 user=postgres dbname=rental sslmode=disable password=$(PASSWORD)" up

migrate-down:
	@goose -dir ./migrates postgres "host=localhost port=5433 user=postgres dbname=rental sslmode=disable password=$(PASSWORD)" down
