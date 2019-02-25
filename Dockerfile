# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/spark-cluster
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/
COPY dashboard/ dashboard/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager spark-cluster/cmd/manager
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o dashboard spark-cluster/dashboard/backend

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /
COPY --from=builder /go/src/spark-cluster/manager .
COPY --from=builder /go/src/spark-cluster/dashboard .
#ENTRYPOINT ["/manager"]
