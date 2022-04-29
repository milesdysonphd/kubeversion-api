FROM golang:1.18-alpine as builder

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY main.go main.go
COPY pkg pkg
RUN go build -o kubeversion-api main.go

FROM golang:1.18-alpine as final
COPY --from=builder kuberversion-api .
CMD ["/kuberversion-api"]