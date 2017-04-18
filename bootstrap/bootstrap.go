package bootstrap

import (
	"log"
	"github.com/grofers/go-codon/bootstrap/golang"
)

type bootstrappable interface {
	Bootstrap() bool
}

var language_map = map[string]bootstrappable {
	"golang": &golang.Bootstrapper,
}

func Bootstrap(lang string) bool {
	bs, ok := language_map[lang]
	if !ok {
		log.Println("Support for language ", lang, " not implemented yet")
	}
	return bs.Bootstrap()
}
