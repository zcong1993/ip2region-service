package main

import (
	"context"
	"fmt"
	"log"

	"github.com/zcong1993/ip2region-service/pb"
	"google.golang.org/grpc"
)

func main() {
	c, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	client := pb.NewIP2RegionClient(c)

	resp, err := client.Search(context.Background(), &pb.IP{Ip: "127.0.0.1"})

	if err != nil {
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("%+v\n", resp)
}
