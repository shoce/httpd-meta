FROM golang:1.26-alpine AS build
RUN mkdir /src/
COPY go.mod *.go /src/
WORKDIR /src/
RUN go build -o httpd-meta

FROM alpine:3
COPY --from=build /src/httpd-meta /bin/httpd-meta
EXPOSE 80
ENTRYPOINT ["/bin/httpd-meta"]
