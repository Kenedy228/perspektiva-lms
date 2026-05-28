FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend

EXPOSE 8080

CMD ["go", "run", "./backend/cmd/lms-api"]
