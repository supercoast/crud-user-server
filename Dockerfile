FROM golang

ENV GOOGLE_APPLICATION_CREDENTIALS="/app/src/go/go-storage-admin-tester.json" \
    GCP_PROJECT_ID="copa-cloud-b8e409187c9d5d11" \
    GCP_BUCKET_NAME="my-test-toptoajdfsljfk"

WORKDIR /app/src/go

COPY go.mod go.sum ./

RUN go mod download

COPY srv-config.yaml /opt

COPY . .

RUN go build -o crud-user-server .

EXPOSE 8080

CMD [ "./crud-user-server" ]