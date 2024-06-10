export PORT=:8080

run:
	templ generate
	PORT="$(PORT)" go run main.go

dev:
	PORT="$(PORT)" templ generate --watch --cmd "go run main.go"
