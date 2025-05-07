FROM golang:1.24.1
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
CMD ["golang"]
EXPOSE 80

FROM alpine:latest
WORKDIR /app
COPY --from=builder /soldatova64 /app/
COPY web ./web/
COPY .env .
EXPOSE 80
CMD ["./soldatova64"]