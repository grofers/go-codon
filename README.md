# go-codon
Generates Go server code from a combination of REST and Workflow DSLs.

[![Build Status](https://travis-ci.org/grofers/go-codon.svg?branch=master)](https://travis-ci.org/grofers/go-codon)

A codon service has three components:
- `Server`: Accepts and validates HTTP requests
- `Clients`: Clients for upstream services which this service consumes
- `Workflows`: Procedures for every REST endpoint of this server which consume Clients and other custom actions.

Server and Client side specifications are written in Swagger. Swagger code generation is done through go-swagger. Workflow is written in `Flow`, a Mistral inspired workflow specification in YAML. Its specification can be found [here](https://github.com/grofers/go-codon/wiki/Workflow-DSL-Specification).

Check out [wiki](https://github.com/grofers/go-codon/wiki) section for more information. Follow [this tutorial](https://github.com/grofers/go-codon/wiki/Codon:-REST-Workflow-Framework) for a very basic example on how to use this tool.

## Projects go-codon would not exist without
(Or just projects I am really thankful for)
- [go-swagger](https://github.com/go-swagger/go-swagger): Provides code generators for client and server side components using [Swagger](https://swagger.io/) specification.
- [go-jmespath](https://github.com/jmespath/go-jmespath): Allows for easy querying and manipulation of json objects in workflows.
- [Pongo2](https://github.com/flosch/pongo2): Django template renderer. Used for templates and workflow expressions in codon.
- [Mistral DSL](https://docs.openstack.org/mistral/latest/): A workflow spec used for infrastructure orchestration. Codon's workflow DSL is inspired from Mistral's but modified for use in REST context.
- [mapstructure](https://github.com/mitchellh/mapstructure)
