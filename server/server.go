package main

import (
	pb "github.com/youtangai/grpcBidirectional/proto"
	"io"
	"net"
	"log"
	"google.golang.org/grpc"
)

type GreetService struct {
	FirstNames []string
}

func (service *GreetService) Greet(stream pb.GreetService_GreetServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("reciev %s\n", in.Message)
		for _, firstname := range service.FirstNames{
			sendMessage := firstname + " " + in.Message
			if err := stream.Send(&pb.Res{Message:sendMessage}); err != nil {
				return err
			}
			log.Printf("send %s\n", sendMessage)
		}
	}
} 

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to start listening 8080: %v", err)
	}
	grpcServer := grpc.NewServer()
	names := []string{"yota", "ryuhei", "miyoshi", "taichi", "hiroto"}
	greetService := &GreetService{
		FirstNames: names,
	}
	pb.RegisterGreetServiceServer(grpcServer, greetService)
	log.Println("start grpc server")
	grpcServer.Serve(lis)
}
