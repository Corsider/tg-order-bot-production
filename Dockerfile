FROM golang:1.20

WORKDIR /app

COPY . .

CMD ["go", "run", "main.go", "structs.go", "bot.go", "dbFuntions.go"]