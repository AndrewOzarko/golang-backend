FROM golang:1.20

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download


RUN go build -o main ./cmd/server
RUN chmod +x ./main

CMD ["./main"]