FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go build -o /appli ./examples/client

EXPOSE 2000

CMD ["/appli"]
