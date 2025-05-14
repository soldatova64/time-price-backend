FROM golang:1.24.1 as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
CMD ["golang"]
EXPOSE 80