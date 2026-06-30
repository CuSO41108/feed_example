FROM golang:1.23-alpine AS build

RUN apk add --no-cache ca-certificates git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./cmd/worker

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /out/worker /app/worker
COPY --from=build /src/docs /app/docs
COPY --from=build /src/web /app/web

EXPOSE 8080

CMD ["/app/api"]
