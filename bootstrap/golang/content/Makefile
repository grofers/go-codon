installdependencies:
	./scripts/dependencies.sh

ibuild:
	go build -o {{pascalize .ProjectName}} server/cmd/{{.ProjectName}}-service-server/main.go

igenerate:
	./scripts/generate.sh

generate: igenerate

build: installdependencies ibuild

.PHONY: installdependencies ibuild igenerate generate build
