FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod .
COPY main.go .
RUN go build -o /app

FROM alpine:3.20
COPY --from=build /app /app
ENTRYPOINT ["/app"]
