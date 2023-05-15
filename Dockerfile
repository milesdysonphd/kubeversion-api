FROM golang:1.20-alpine as builder

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY pkg/ pkg/
COPY controllers/ controllers/

RUN GOOS=linux GOARCH=amd64 go build -o kubeversion-api main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /workspace/kubeversion-api .
CMD ["/app/kubeversion-api"]