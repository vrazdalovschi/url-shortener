FROM golang:1.14-alpine AS build
RUN mkdir /go/src/url-shortener
ADD ./ /go/src/url-shortener
WORKDIR /go/src/url-shortener
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main ./main.go

FROM scratch as runtime
COPY --from=build /main /
EXPOSE 8080
ENTRYPOINT ["/main"]
