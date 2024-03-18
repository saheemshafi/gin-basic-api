FROM golang

WORKDIR /app
ENV GIN_MODE=release

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]