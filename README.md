# gRPC-Go-Sample

To test the server using evans cli :-
`evans -p 50051 -r`


To generate TLS certificates, run `sh generateCertificates.sh`


No. of RPC per second is coming about 2000 after doing the perf testing

To build a executable using a preferred name use `-o`,

`go build -o greet-server server/server.go`

To dockerize server :-
run `docker build -t greet-server -f Dockerfile.server .`
`docker run -d -p 50051:50051 -v /Users/sumisaha/go/src/github.com/saha/grpc-go-course/greet/logs/server:/app/logs greet-server`
