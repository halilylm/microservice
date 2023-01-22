# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-X 'main.release=`git rev-parse --short=8 HEAD`'" -o /bin/server ./app

FROM scratch

WORKDIR /

COPY --from=build /bin/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]