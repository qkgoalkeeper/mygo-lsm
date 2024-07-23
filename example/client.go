package main

import (
	"context"
	"fmt"
	"log"
	_ "time"

	pb "github.com/whuanle/lsm/proto"

	"google.golang.org/grpc"
)

func main() {
	// 建立与gRPC服务器的连接
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewKeyValueServiceClient(conn)

	// 调用Set方法
	setRequest := &pb.SetRequest{
		Key: "example_key",
		Value: &pb.TestValue{
			A: 1,
			B: 2,
			C: 3,
			D: "example_value",
		},
	}
	setResponse, err := client.Set(context.Background(), setRequest)
	if err != nil {
		log.Fatalf("could not set value: %v", err)
	}
	fmt.Printf("Set Response: %v\n", setResponse)

	// 调用Search方法
	searchRequest := &pb.SearchRequest{
		Key: "example_key",
	}
	searchResponse, err := client.Search(context.Background(), searchRequest)
	if err != nil {
		log.Fatalf("could not search value: %v", err)
	}
	fmt.Printf("Search Response: %v\n", searchResponse)
}
