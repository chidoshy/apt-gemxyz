GO_CMD_MAIN = main.go

migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-force:
	echo \# make migrate-force version="$(name)"
	go run $(GO_CMD_MAIN) migrate force "$(version)"

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

up:
	docker-compose up -d

down:
	docker-compose down
