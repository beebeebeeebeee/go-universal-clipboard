FROM golang:1.22-bullseye as builder

ENV GO111MODULE=on

RUN mkdir -p /microservice
ADD . /microservice
WORKDIR /microservice

RUN go mod download
COPY . .
RUN go build ./cmd/app/main.go

FROM ubuntu:22.04
RUN mkdir -p /microservice/internal/app/infrastructure/gin/static
WORKDIR /microservice
COPY --from=builder /microservice/main .
COPY --from=builder /microservice/config.yaml .
COPY --from=builder /microservice/internal/app/infrastructure/gin/static /microservice/internal/app/infrastructure/gin/static
ENTRYPOINT ["./main"]
