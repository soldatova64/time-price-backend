FROM golang:1.24.1
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
CMD ["golang"]
EXPOSE 80