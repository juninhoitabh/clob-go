mock:
	mockgen -source=$(INTERFACE) -destination=$(DEST) -package=mocks
	# INTERFACE ?= internal/domain/account/repository.go
	# DEST ?= internal/domain/account/mocks/account_repository.mock.go

# SWAGGER
swagGenerate: 
	swag init -g internal/infra/http-server/router/router.go

# RUN Tests and Lint
test-lint:
	./gosweep.sh