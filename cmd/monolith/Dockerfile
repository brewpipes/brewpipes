FROM golang:1.25.5

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go run ./cmd/monolith

EXPOSE 8080

ENTRYPOINT ["./brewpipes"]
