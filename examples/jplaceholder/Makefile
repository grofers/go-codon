installdependencies:
	./scripts/dependencies.sh

ibuild:
	go build -o Jplaceholder server/cmd/jplaceholder-service-server/main.go

igenerate:
	./scripts/generate.sh

generate: installdependencies igenerate

build: installdependencies igenerate ibuild

.PHONY: installdependencies ibuild igenerate generate build
