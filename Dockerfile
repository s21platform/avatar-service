FROM golang:1.22 as builder

#RUN apt-get update && apt-get install -y libwebp-dev

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go

FROM alpine

RUN apk add --no-cache libwebp-dev

WORKDIR /app

COPY --from=builder /usr/src/service/build/main /app
RUN apk add --no-cache gcompat

CMD ["./main"]