FROM public.ecr.aws/zomato/golang:1.19-bullseye

ENV GOBIN $GOPATH/bin
WORKDIR $GOPATH/src/github.com/grofers/go-codon
COPY . .

RUN make install
