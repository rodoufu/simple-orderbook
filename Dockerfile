FROM golang:1.19.1-buster as builder
WORKDIR /app
ADD . .
RUN go test -v --cover ./...
RUN cd cmd/simplebook && go build

FROM ubuntu:22.04 as runner
WORKDIR /app
COPY --from=builder /app/cmd/simplebook/simplebook .
ENTRYPOINT ["./simplebook"]
