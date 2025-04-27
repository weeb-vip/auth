generate: gql mocks
gql:
	go run github.com/99designs/gqlgen generate
create-migration:
	go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir internal/migrations/scripts $(name)

migrate:
	go run cmd/cli/main.go db migrate

mocks:
	go get go.uber.org/mock/mockgen/model
	go install go.uber.org/mock/mockgen@latest
	mockgen -destination=./mocks/mock_credentials.go -package=mocks github.com/weeb-vip/auth/internal/services/credential Credential
	mockgen -destination=./mocks/mock_tokenizer.go -package=mocks github.com/weeb-vip/auth/internal/jwt Tokenizer
	mockgen -destination=./mocks/mock_user_client.go -package=mocks github.com/weeb-vip/auth/internal/services/user_client UserClientInterface
