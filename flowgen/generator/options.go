package generator

import (
	"os"
)

type GenOpts struct {
	Spec			string
	Dest			string
	Templates		string
}

func (o *GenOpts) EnsureDefaults() {
	if o.Spec == "" {
		if _, err := os.Stat("flow.yml"); err == nil {
			o.Spec = "flow.yml"
		} else {
			o.Spec = "flow.yaml"
		}
	}
	if o.Dest == "" {
		o.Dest = "flow.go"
	}
}

func (o *GenOpts) Verify() error {
	if _, err := os.Stat(o.Spec); os.IsNotExist(err) {
		return err
	}

	return nil
}
