package generator

import (
	"github.com/grofers/go-codon/flowgen/shared"
)

func Process(opts *GenOpts) error {
	opts.EnsureDefaults()
	if err := opts.Verify(); err != nil {
		return err
	}

	spec, err := shared.ReadSpec(opts.Spec)
	if err != nil {
		return err
	}

	spec_ptr := &spec

	post_spec, err := spec_ptr.Process()
	if err != nil {
		return err
	}

	err = Generate("go", opts, &post_spec)
	if err != nil {
		return err
	}

	return nil
}
