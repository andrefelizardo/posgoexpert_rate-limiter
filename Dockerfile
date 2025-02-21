FROM golang:1.23.5
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
CMD ["go", "test", "./..."]
