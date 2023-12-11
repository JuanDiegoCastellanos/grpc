package main

import (
	"context"
	"github.com/JuanDiegoCastellanos/grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func main() {

	// Conectando al server
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()
	// Creando cliente
	c := testpb.NewTestServiceClient(cc)

	//DoUnary(c)
	//DoClientStreamin(c)
	//DoServerStreaming(c)
	DoBidirectionalStreaming(c)
}

// DoUnary Funcion para cliente Unario
func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "test1",
	}
	// Esto es grpc como tal, ya que llama a metodos que parecieran declarados en este scope
	// pero que en realidad estan a nivel de server
	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest: %v", err)
	}
	log.Printf("response from server: %v", res)
}

// Streaming del lado del cliente
// DoClientStreaming
func DoClientStreamin(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q7",
			Answer:   "ZaiDich",
			Question: "Rammstein",
			TestId:   "test3",
		},
		{
			Id:       "q8",
			Answer:   "Azure",
			Question: "rammstein",
			TestId:   "test3",
		},
		{
			Id:       "q9",
			Answer:   "Frontend",
			Question: "Que hacen los gays?",
			TestId:   "test3",
		},
	}
	stream, err := c.SetQuestion(context.Background())
	if err != nil {
		log.Fatalf("error while Calling SetQuestions: %v", err)
	}
	for _, question := range questions {
		log.Printf("Sending a question: %v ...", question.Id)
		err := stream.Send(question)
		if err != nil {
			log.Fatalf("error while sending a question: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while closing streaming: %v", err)
	}
	log.Printf("Server response: %v", msg)
}

// Streaming desde el server
func DoServerStreaming(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "test3",
	}
	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetStudentsPerTest: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading the stream: %v", err)
		}
		log.Printf("response from server: %v", msg)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer:    "70",
		TestId:    "test3",
		StudentId: "student1",
	}
	numbeOfQuestions := 4

	// Channel de control, sirve para bloquear las goroutines
	waitChanel := make(chan struct{})
	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("error while calling TakeTest: %v", err)
	}
	go func() {
		for i := 0; i < numbeOfQuestions; i++ {
			err := stream.Send(&answer)
			if err != nil {
				log.Fatalf("error while sending answer: %v", err)
			}
			time.Sleep(2 * time.Second)
		}
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error reading streaming: %v", err)
			}
			log.Printf("response from server: %v", res)
		}
		close(waitChanel)
	}()
	// bloquear el programa, se va a quedar esperando que termine el servidor de hacer streaming
	<-waitChanel
}
func DoBidireccionalStreamingv2(c testpb.TestServiceClient) {
	startTest := &testpb.TakeTestRequest{
		TestId: "t1",
	}
	testAnswer := &testpb.TakeTestRequest{
		TestId: "t1",
		Answer: "asdasdas",
	}
	stream, err := c.TakeTest(context.Background())

	if err != nil {
		log.Fatalf("Error while calling TakeTest: %v", err)
	}

	stream.Send(startTest)

	for {
		msg, err := stream.Recv()

		// Last message from server
		if err == io.EOF || msg.GetOk() {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("question %v", msg)
		stream.Send(testAnswer)
	}

	closeErr := stream.CloseSend()

	if closeErr != nil {
		log.Fatalf("Error while closing stream: %v", closeErr)
	}

	log.Println("Connection closed")
}
