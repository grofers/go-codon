package generator

import (
	"errors"
	"github.com/grofers/go-codon/flowgen/languages"
	"github.com/grofers/go-codon/flowgen/shared"
)

func Generate(language string, opts *GenOpts, post_spec *shared.PostSpec) error {
	var err error
	var gen languages.Generator
	switch language {
	case "go":
		gen = &languages.GoGenerator {
			Data: post_spec,
			Dest: opts.Dest,
			Templates: opts.Templates,
		}
	default:
		gen = nil
	}
	if gen == nil {
		err = errors.New("Language not implemented yet")
	} else {
		err = gen.Generate()
	}
	return err
}
