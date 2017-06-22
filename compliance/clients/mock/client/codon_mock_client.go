package client

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/grofers/go-codon/testing/clients/mock/client/operations"
)

func NewHTTPClientWithConfigMap(formats strfmt.Registry, cfgmap *map[string]string) *CodonMock {
	return New()
}

func New() *CodonMock {
	cli := new(CodonMock)
	cli.Operations = operations.New()
	return cli
}

type CodonMock struct {
	Operations *operations.Client
}
