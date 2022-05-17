# go-codon
Generates Go Server Code from a combination of REST and Workflow DSLs.

[![Build Status](https://travis-ci.org/grofers/go-codon.svg?branch=master)](https://travis-ci.org/grofers/go-codon)

A codon service has three components:
- `Server`: Accepts and validates HTTP requests
- `Clients`: Clients for upstream services which this service consumes
- `Workflows`: Procedures for every REST endpoint of this server which consume Clients and other custom actions.

Server and Client side specifications are written in Swagger. Swagger code generation is done through go-swagger. Workflow is written in `Flow`, a Mistral inspired workflow specification in YAML. Its specification can be found [here](https://github.com/grofers/go-codon/wiki/Workflow-DSL-Specification).

Check out [wiki](https://github.com/grofers/go-codon/wiki) section for more information. Follow [this tutorial](https://github.com/grofers/go-codon/wiki/Codon:-REST-Workflow-Framework) for a very basic example on how to use this tool.

## Installation
Set up your Golang development environment ([Getting Started](https://golang.org/doc/install)). Set your `GOPATH` and `GOBIN` directories. Also add `GOBIN` to your `PATH` so that golang tools can be used in command line.

Download the latest binary from Github releases and put it in your `GOBIN` directory. Or to install from source do:
```sh
mkdir -p $GOPATH/src/github.com/grofers
cd $GOPATH/src/github.com/grofers
git clone git@github.com:grofers/go-codon.git
cd go-codon
make install
```

## Example
This is what a workflow looks like (for an API to get posts and the comments for each post concurrently):
```yaml
name: get_posts_comments
start:
    - get_posts
tasks:
    get_posts:
        action: clients.jplaceholder.get_posts
        input:
            userId: <%jmes main.userId %>
        publish:
            posts: <%jmes action %>
        on-success:
            - get_all_comments: true
    get_comments:
        action: clients.jplaceholder.get_comments
        input:
            postId: <%jmes main.postId %>
        publish:
            comments: <%jmes action %>
    get_all_comments:
        with-items: <%jmes main.posts %>
        loop:
            task: get_comments
            input:
                postId: <%jmes item.id %>
            publish:
                combined: <%jmes {"post_details":item,"comments":task.comments} %>
output:
    body: <%jmes main.combined %>
    status_code: 200
```
To run this example checkout [examples](https://github.com/grofers/codon-examples).

## Projects go-codon would not exist without
(Or just projects I am really thankful for)
- [go-swagger](https://github.com/go-swagger/go-swagger): Provides code generators for client and server side components using [Swagger](https://swagger.io/) specification.
- [go-jmespath](https://github.com/jmespath/go-jmespath): Allows for easy querying and manipulation of json objects in workflows.
- [Pongo2](https://github.com/flosch/pongo2): Django template renderer. Used for templates and workflow expressions in codon.
- [Mistral DSL](https://docs.openstack.org/mistral/latest/): A workflow spec used for infrastructure orchestration. Codon's workflow DSL is inspired from Mistral's but modified for use in REST context.
- [mapstructure](https://github.com/mitchellh/mapstructure)
