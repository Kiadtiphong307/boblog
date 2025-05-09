FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# ติดตั้ง Air
RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 8080
CMD ["air"]