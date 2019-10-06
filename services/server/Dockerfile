FROM golang:1.12.9-buster AS builder
ADD . /app
WORKDIR /app
RUN go get -d -v github.com/gin-gonic/gin \
	github.com/spf13/viper \
	github.com/xanzy/go-gitlab
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 8000