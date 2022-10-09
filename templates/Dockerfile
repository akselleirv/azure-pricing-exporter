FROM golang:1.19 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN  CGO_ENABLED=0 GOOS=linux go build -a -o app .

FROM gcr.io/distroless/static

EXPOSE 8080

COPY --from=builder /app/app .

ENTRYPOINT ["./app"]