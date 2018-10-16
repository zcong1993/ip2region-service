package main

import (
	"context"
	"github.com/zcong1993/ip2region-service/service"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/zcong1993/ip2region-service/pb"
	"google.golang.org/grpc"
)

func runGatewayServer(rpcPort, port string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterIP2RegionHandlerFromEndpoint(ctx, mux, rpcPort, opts)

	if err != nil {
		log.Fatal("Serve http error:", err)
	}

	http.ListenAndServe(port, mux)
}

func runRpcServer(port string) {
	ss := grpc.NewServer()
	pb.RegisterIP2RegionServer(ss, service.NewIP2RegionService(os.Getenv("DB_PATH")))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	if err = ss.Serve(listener); err != nil {
		log.Fatal("ListenTCP error:", err)
	}
}

func main() {
	if os.Getenv("GATEWAY") == "true" {
		println("run gateway server on :9009")
		go runGatewayServer(":1234", ":9009")
	}
	runRpcServer(":1234")
}
