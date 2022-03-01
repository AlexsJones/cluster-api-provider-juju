FROM alpine:latest

RUN apk add --no-cache git make musl-dev go curl

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin
RUN curl -LO https://launchpad.net/juju/2.9/2.9.0/+download/juju-2.9.0-linux-amd64.tar.xz
RUN tar xf juju-2.9.0-linux-amd64.tar.xz
RUN install -o root -g root -m 0755 juju /usr/local/bin/juju
WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

ENTRYPOINT ["/workspace/manager"]
