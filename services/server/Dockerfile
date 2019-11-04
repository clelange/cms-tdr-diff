FROM golang:1.13.4-buster AS builder
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 8000