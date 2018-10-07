# -----------------
# Builder Container
# -----------------

FROM golang:1.8 AS builder

ENV GOBIN $GOPATH/bin

RUN mkdir -p $GOPATH/src/github.com/grofers
COPY . $GOPATH/src/github.com/grofers/go-codon/

RUN apt-get update
RUN apt-get install unzip
RUN wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.zip
RUN unzip glide-v0.12.3-linux-amd64.zip
RUN mv linux-amd64/glide /usr/local/bin/

WORKDIR $GOPATH/src/github.com/grofers/go-codon/

RUN make install

# ---------------
# Final Container
# ---------------

FROM golang:1.8

WORKDIR /go

COPY --from=builder /go/bin/codon /bin/codon
COPY --from=builder /go/bin/go-bindata /bin/go-bindata
