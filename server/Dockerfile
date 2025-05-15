FROM golang:1.24.3-alpine as development

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum .

RUN go mod download

CMD ["air"]
