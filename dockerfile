FROM golang:1.22.2

WORKDIR /app

RUN apt update && apt install -y sqlite3

COPY go.sum .
COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
