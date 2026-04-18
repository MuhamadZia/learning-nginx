FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o app

EXPOSE 8080

CMD [ "./app" ]

# run healthcheck for awarness nginx
HEALTHCHECK --interval=10s --timeout=2s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1