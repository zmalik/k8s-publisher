FROM golang:1.9.2-alpine as builder
WORKDIR /go/src/github.com/zmalik/k8s-publisher/
ADD . .
RUN go build -o k8s-publisher

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/zmalik/k8s-publisher/k8s-publisher .
CMD [ "./k8s-publisher" ]