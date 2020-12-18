# Multi-stage build step 0
FROM golang:latest AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /server

# Final stage
FROM debian:buster
EXPOSE 8085
WORKDIR /
COPY --from=build-env /server /
CMD ["/server"]
