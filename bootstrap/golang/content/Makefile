installdependencies:
	./scripts/dependencies.sh

ibuild:
	go build -o {{pascalize .ProjectName}} server/cmd/{{.ProjectName}}-service-server/main.go

igenerate:
	codon generate

generate: installdependencies igenerate

build: installdependencies igenerate ibuild

.PHONY: installdependencies ibuild igenerate generate build
