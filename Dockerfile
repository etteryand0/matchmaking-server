FROM golang:1.22-alpine

WORKDIR /server

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV PORT=8000

RUN go build -v -o ./server

CMD ["./server"]