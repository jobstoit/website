FROM golang:1.17-alpine AS builder

ADD . /build

WORKDIR /build

RUN mkdir -p /out && \
	go mod tidy && \
	go build -o /out/website /build/cmd/website

FROM alpine:latest

ENV WEBSITE_PORT=8080
ENV WEBSITE_DATABASE_CONNECTION_STRING=postgres://user:password@localhost:5432/dbname?sslmode=disable
ENV WEBSITE_OID_PROVIDER=""
ENV WEBSITE_OID_CLIENT_ID=""
ENV WEBSITE_OID_CLIENT_SECRET=""
ENV WEBSITE_OID_STATE_STRING="randomChars"

COPY --from=builder /out/website /usr/local/bin/website

CMD ["/usr/local/bin/website"]
