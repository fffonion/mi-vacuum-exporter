FROM golang:1.14-buster

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch

COPY --from=0 /app/mi-vacuum-exporter /mi-vacuum-exporter

EXPOSE 9234

ENTRYPOINT ["/mi-vacuum-exporter"]