FROM golang:1.18
WORKDIR /testgrpc
COPY . .
RUN go mod download
RUN go build -o testtask2 ./mainpage.go
CMD ["./testtask2"]