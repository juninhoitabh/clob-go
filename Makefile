mock:
	mockgen -source=$(INTERFACE) -destination=$(DEST) -package=mocks
	# INTERFACE ?= internal/domain/account/repository.go
	# DEST ?= internal/domain/account/mocks/account_repository.mock.go