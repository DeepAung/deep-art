DATABASE_URL = sqlite://$(abspath ./db.db)
MIGRATION_URL = file://$(abspath ./pkg/db/migrations)

air:
	air -c .air.toml
tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
templ:
	templ generate --watch --proxy="http://localhost:3000"
tidy:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	templ generate
	go mod tidy

migrate.up:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose up
migrate.down:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose down
migrate.goto:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose goto $(VERSION)

jet:
	jet -source=sqlite -dsn="./db.db" -schema=dvds -path=./.gen
