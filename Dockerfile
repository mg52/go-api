FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o  /go-api

EXPOSE 3000

CMD [ "/go-api" ]