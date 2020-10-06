FROM golang:1.14-alpine AS build
WORKDIR /go/src/url-shortener
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/url-shortener ./main.go

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /go
COPY --from=build /go/src/url-shortener/bin /go/bin
COPY --from=build /go/src/url-shortener/api /go/api
EXPOSE 8080
ENTRYPOINT /go/bin/url-shortener