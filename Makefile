goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	@read -p "Name for the change (e.g. add_column): " name; \
	goose -dir db/migrations/ create $${name:-<name>} sql

sqlc:
	sqlc generate

mocks:
	go install github.com/golang/mock/mockgen@v1.6.0
	GO111MODULE=on mockgen -source infraestructure/db/repos/querier.go -destination test/mocks/querier.go -package mocks
	GO111MODULE=on mockgen -source domain/services/auth/authorization.go -destination test/mocks/authorization.go -package mocks
	GO111MODULE=on mockgen -source domain/services/transaction/transaction.go -destination test/mocks/transaction.go -package mocks
	GO111MODULE=on mockgen -source domain/clients/bank.go -destination test/mocks/bank.go -package mocks

docker-up-np:
	docker-compose -f docker-compose.yml up -d --build --wait --scale payment=0
	@echo "payment api up and running"

docker-up:
	docker-compose -f docker-compose.yml up -d --build --wait
	@echo "payment api up and running"
	@echo "Swagger" http://localhost:8090/swagger/index.html

docker-down:
	docker-compose -f docker-compose.yml down --remove-orphans --volumes
	@echo "payment down"

docs:
	go install github.com/swaggo/swag/cmd/swag
	swag init --parseVendor -g cmd/main.go --output api/docs
