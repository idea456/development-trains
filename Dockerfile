FROM golang:1.23-bullseye AS build

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o development-trains ./cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /app/development-trains development-trains

ENTRYPOINT ["./development-trains"]