export PORT=:8080
export GO_OPTIONS="-tags=viper_bind_struct"

run:
	templ generate
	PORT="$(PORT)" go run ${GO_OPTIONS} main.go

dev:
	PORT="$(PORT)" templ generate --watch --cmd "go run ${GO_OPTIONS} main.go"

build:
	go build ${GO_OPTIONS} -ldflags="-X 'main.Version=$VERSION' -X 'main.Commit=$COMMIT' -X 'main.BuildTime=$BUILD_TIME'" -o bin/golink .
