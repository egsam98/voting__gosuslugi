DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_HOST ?= db
DB_NAME ?= postgres

ifneq (,$(wildcard .env))
	include .env
	export
endif

db-create-migration: ## Create migration file in db/migrations directory. Migration should be named by "name" argument. Example: create-migration name=create_foos
	docker run -v "${PWD}/db/migrations:/migrations" \
		--network host migrate/migrate \
		-path=/migrations \
 		create -ext \
 		sql -dir \
 		/migrations $(name)
	sudo chmod 766 -R db/migrations

db-migrate:
	docker run -v "${PWD}/db/migrations:/migrations" \
		--network host \
		migrate/migrate \
        -path=/migrations/ \
        -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable \
        up

db-rollback:
	docker run -v "${PWD}/db/migrations:/migrations" \
		--network host migrate/migrate \
		-path=/migrations \
		-database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable \
		down 1

db-sqlc:
	sqlc generate
