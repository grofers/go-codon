package main

import "github.com/grofers/go-codon/cmd"

// All the generate directives here
//go:generate go-bindata -prefix bootstrap/golang/content/ -pkg golang -o bootstrap/golang/bindata.go bootstrap/golang/content/...
//go:generate go-bindata -prefix generator/golang/content/ -pkg golang -o generator/golang/bindata.go generator/golang/content/...

func main() {
	cmd.Execute()
}
