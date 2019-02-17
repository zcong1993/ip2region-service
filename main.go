package main

import (
	"context"
	"fmt"
	"github.com/zcong1993/ip2region-service/service"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/zcong1993/ip2region-service/pb"
	"google.golang.org/grpc"
)

func initJaegerTracing() {
	svcAddr := os.Getenv("JAEGER_SERVICE_ADDR")
	if svcAddr == "" {
		fmt.Println("jaeger initialization disabled.")
		return
	}

	view.RegisterExporter(&exporter.PrintExporter{})

	// Register the views to collect server request count.
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: fmt.Sprintf("http://%s", svcAddr),
		Process: jaeger.Process{
			ServiceName: "ip2region-service",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	fmt.Println("jaeger initialization completed.")
}

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
	ss := grpc.NewServer(grpc.StatsHandler(&ocgrpc.ServerHandler{}))
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
	go initJaegerTracing()
	if os.Getenv("GATEWAY") == "true" {
		println("run gateway server on :9009")
		go runGatewayServer(":1234", ":9009")
	}
	runRpcServer(":1234")
}
