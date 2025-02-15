FROM golang:1.22.2

WORKDIR /app

RUN apt update && apt install -y sqlite3

COPY go.mod .
COPY main.go .
COPY db.sql .

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
