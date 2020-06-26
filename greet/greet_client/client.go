package main

import (
	"fmt"
	"go-grpc-course/greet/greetpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Starting greet client!!")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//
	//doServerServerSstreaming(c)
	//
	//doClientSstreaming(c)
	doBiDiStreaming(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doClientSstreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mohit",
				LastName:  "Joshi",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
				LastName:  "",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "Joshi",
			},
		}, &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
				LastName:  "",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		panic(err)
	}

	for _, req := range requests {
		fmt.Println("Sending stream request:\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	fmt.Println("Long Greet response is ", res)

}

func doServerServerSstreaming(c greetpb.GreetServiceClient) {
	req := greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohit",
			LastName:  "Joshi",
		},
	}
	res, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println("SuccessGreetStream: ", msg)
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Mohit",
			LastName:  "Joshi",
		},
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println("SuccessGreet: ", res.Result)
}
