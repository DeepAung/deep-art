DATABASE_URL = sqlite3://$(abspath ./db.db)
MIGRATION_URL = file://$(abspath ./migrations)

air:
	air -c .air.toml
tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch
tailwind.reset:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css
templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false

migrate.goto:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose goto $(VERSION)
migrate.up:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose up
migrate.down:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose down
migrate.force:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) -verbose force $(VERSION)
migrate.version:
	migrate -database $(DATABASE_URL) -source $(MIGRATION_URL) version
migrate.reset:
	make migrate.down && make migrate.up

jet:
	jet -source=sqlite -dsn="./db.db" -schema=dvds -path=./.gen
jet.test:
	jet -source=sqlite -dsn="./test.db" -schema=dvds -path=./.gen

tidy:
	make migrate.reset
	make jet
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify
	templ generate
	go mod tidy
docker.build: tidy
	docker build --pull -t deep-art:latest .
docker.run:
	docker run --env-file .env.prod -p 3000:3000 --name deep-art deep-art:latest
docker.push:
	docker tag deep-art:latest $(IMAGE_URL)
	docker push $(IMAGE_URL)
