air:
	air

# DATABASE_URL = postgres://myuser:mypassword@0.0.0.0:5432/mydb?sslmode=disable
# MIGRATION_FOLDER = file://$(abspath ./pkg/databases/migrations)
#
# run: db.start fiber.dev
#
# fiber.dev:
# 	air -c .air.dev.toml
# fiber.prod:
# 	air -c .air.prod.toml
#
# db.create:
# 	docker run \
# 		--name deep_art_db_dev \
# 		-e POSTGRES_USER=myuser \
# 		-e POSTGRES_PASSWORD=mypassword \
# 		-e POSTGRES_DB=mydb \
# 		-p 5432:5432 \
# 		-d postgres:alpine
# db.start:
# 	docker start deep_art_db_dev
# db.stop:
# 	docker stop deep_art_db_dev
#
# migrate.up:
# 	migrate -database $(DATABASE_URL) -source $(MIGRATION_FOLDER) -verbose up
# migrate.down:
# 	migrate -database $(DATABASE_URL) -source $(MIGRATION_FOLDER) -verbose down
# migrate.goto:
# 	migrate -database $(DATABASE_URL) -source $(MIGRATION_FOLDER) -verbose goto $(version)
