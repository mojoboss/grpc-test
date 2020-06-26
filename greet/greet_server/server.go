package main

import (
	"fmt"
	"go-grpc-course/greet/greetpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function was invoked with: ", req)
	firstName := req.GetGreeting().FirstName
	result := "Hello " + firstName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("Greet Everyone function was invoked: ")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		firstName := req.GetGreeting().FirstName
		result := "Hello " + firstName + "!\n"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("Greet Many function was invoked with: ", req)
	firstName := req.GetGreeting().FirstName
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName
		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(&res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("Greet Many function was invoked with streaming request ")
	result := "Hello  "
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// FINISHED READING CLIENT STREAM
			stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
			break
		}
		if err != nil {
			return err
		}
		firstName := req.Greeting.FirstName
		result += firstName + "!"
	}

	return nil
}

func main() {
	fmt.Println("Starting greet service!!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
