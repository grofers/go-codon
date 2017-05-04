# go-codon
Generates Go server code from a combination of REST and Workflow DSLs.

A codon service has three components:
- `Server`: Accepts and validates HTTP requests
- `Clients`: Clients for upstream services which this service consumes
- `Workflows`: Procedures for every REST endpoint of this server which consume Clients and other custom actions.

Server and Client side specifications are written in Swagger. Swagger code generation is done through go-swagger. Workflow is written in `Flow`, a Mistral inspired workflow specification in YAML. Its specification can be found here.

Follow [this tutorial](https://github.com/grofers/go-codon/wiki/Codon:-REST-Workflow-Framework) for a very basic example on how to use this tool.
