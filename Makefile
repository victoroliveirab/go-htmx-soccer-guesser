build:
	go build -o tmp/main .

run: build
	APP_ENV=dev ./tmp/main

tailwindcss:
	bun run tailwindcss --config tailwind.config.js -i base.css -o static/css/styles.css --watch

bootstrap:
	python3 cmd/bootstrap_db.py
	python3 cmd/bootstrap_guesses.py

points:
	go run ./cmd/generate_outcomes.go
	go run ./cmd/update_group_ranking.go

test:
	go test -v -failfast ./...
