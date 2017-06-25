package generator

import (
	"log"
	"github.com/grofers/go-codon/generator/shared"
	"github.com/grofers/go-codon/generator/golang"
)

type GenOpts shared.GenOpts

type generatable interface {
	Generate(shared.GenOpts) bool
}

var language_map = map[string]generatable {
	"golang": &golang.Generator,
}

func Generate(opts GenOpts) bool {
	bs, ok := language_map[opts.Language]
	if !ok {
		log.Println("Support for language ", opts.Language, " not implemented yet")
	}
	return bs.Generate(shared.GenOpts(opts))
}
