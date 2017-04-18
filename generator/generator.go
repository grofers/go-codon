package generator

import (
	"log"
	"github.com/grofers/go-codon/generator/golang"
)

type generatable interface {
	Generate() bool
}

var language_map = map[string]generatable {
	"golang": &golang.Generator,
}

func Generate(lang string) bool {
	bs, ok := language_map[lang]
	if !ok {
		log.Println("Support for language ", lang, " not implemented yet")
	}
	return bs.Generate()
}
