package main

import (
	pb "github.com/youtangai/grpcBidirectional/proto"
	"golang.org/x/net/context"
	"time"
	"io"
	"log"
	"google.golang.org/grpc"
)

func runGreet(client pb.GreetServiceClient) {
	lastNames := []string{"nagai", "nanaumi", "miyoshi", "ito", "ito"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.Greet(ctx)
	if err != nil {
		log.Fatalf("failed to greet: %v\n", err)
	}
	waitchan := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive: %v", err)
			}
			log.Printf("recieve! Hello %s\n", in.Message)
		}
	}()
	for _, name := range lastNames {
		if err := stream.Send(&pb.Req{Message: name}); err != nil {
			log.Fatalf("failed to send: %v\n", err)
		}
		log.Printf("send %s\n", name)
	}
	stream.CloseSend()
	<-waitchan
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial %v\n", err)
	}
	defer conn.Close()
	client := pb.NewGreetServiceClient(conn)
	runGreet(client)
}
