FROM golang

ENV GOBIN $GOPATH/bin
WORKDIR $GOPATH/src/github.com/grofers/go-codon
COPY . .

RUN make install
