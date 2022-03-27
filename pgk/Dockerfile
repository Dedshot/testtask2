FROM golang:1.18
WORKDIR /testgrpc
ADD . ./pgk
ADD . ./testtask2/pgk/api
COPY . ./pgk
COPY . ./testtask2/pgk/api

RUN apt-get update
RUN apt-get install dos2unix
RUN go mod init testtask2
RUN go get google.golang.org/api/googleapi/transport
RUN go get google.golang.org/api/youtube/v3
RUN go get google.golang.org/grpc
RUN go build -o grpc ./pgk/grpcserver.go
CMD ["./grpc"]