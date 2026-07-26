FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o main


FROM scratch

COPY --from=builder /app/main /
COPY interface.txt /

ENV TZ=Asia/Shanghai

EXPOSE 6688

ENTRYPOINT ["/main"]