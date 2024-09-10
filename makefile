swag_init:
	swag init -g api/router.go -o api/docs

migration-up:
	migrate -database 'postgres://shahzod:1@localhost:5432/food?sslmode=disable' -path migrations up;

migration-down:
	migrate -path ./migrations/postgres -database 'postgres://shahzod:1@0.0.0.0:5432/food?sslmode=disable' down