package clients

import (
	mock_client "github.com/grofers/go-codon/clients/mock/client"
)


var Mock = mock_client.NewHTTPClientWithConfigMap(nil, nil).Operations
