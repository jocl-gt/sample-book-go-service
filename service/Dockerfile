FROM golang:1.21.1

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 5001

CMD ["./main"]
