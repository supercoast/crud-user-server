FROM golang

WORKDIR /app/src/go

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o crud-user-server .

EXPOSE 8080

CMD [ "./crud-user-server" ]