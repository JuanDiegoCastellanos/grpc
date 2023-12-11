package server

import (
	"context"
	"fmt"
	"github.com/JuanDiegoCastellanos/grpc/studentpb"
	"io"
	"log"
	"time"

	"github.com/JuanDiegoCastellanos/grpc/models"
	"github.com/JuanDiegoCastellanos/grpc/repository"
	"github.com/JuanDiegoCastellanos/grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{
		repo: repo,
	}
}
func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}
	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil

}

func (s *TestServer) SetQuestion(stream testpb.TestService_SetQuestionServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return nil
		}
		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.GetAnswer(),
			Question: msg.GetQuestion(),
			TestId:   msg.GetTestId(),
		}
		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}
func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		//Esto es para verificar si ha sido cerrado o se termino de recibir el stream
		// End of File = El server debe responder porque se termino la conexion
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}
		enrollment := &models.Enrollment{
			StudentId: msg.GetStudentId(),
			TestId:    msg.GetTestId(),
		}
		err = s.repo.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}
func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}
	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}
		err := stream.Send(student)
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Println(err)
			return err
		}

		questions, err := s.repo.GetQuestionPerTest(context.Background(), msg.GetTestId())

		if err != nil {
			return err
		}

		var currentQuestion = &models.Question{}
		i := 0
		lenQuestions := len(questions)
		lenQuestions32 := int32(lenQuestions)

		for {
			if i < lenQuestions {
				currentQuestion = questions[i]
				questionToSend := &testpb.QuestionPerTest{
					Id:       currentQuestion.Id,
					Question: currentQuestion.Question,
					Ok:       false,
					Current:  int32(i + 1),
					Total:    lenQuestions32,
				}

				// Enviamos la question
				err := stream.Send(questionToSend)

				if err != nil {
					return err
				}
				// Aumentamos contador para pasar a la siguiente
				i++
				// Ahora recibimos la respuesta
				answer, err := stream.Recv()

				if err == io.EOF {
					return nil
				}

				if err != nil {
					return err
				}

				log.Println("Answer: ", answer.GetAnswer())
				// Lo siguiente es para guardar en base de datos la respuesta
				answerModel := &models.Answer{
					TestId:     msg.GetTestId(),
					QuestionId: currentQuestion.Id,
					StudentId:  msg.GetStudentId(),
					Answer:     answer.Answer,
					Correct:    answer.Answer == currentQuestion.Answer,
				}

				err = s.repo.SetAnswer(context.Background(), answerModel)

				if err != nil {
					fmt.Println(err)
					return err
				}
			} else { // En este punto ya se envarion todas las questions
				// Se devuelve una QuestionToSend pero vacio y con el flag ok en true
				// ese flag Ok = true significa que todas las preguntas(questions) ya fueron enviadas
				questionToSend := &testpb.QuestionPerTest{
					Id:       "",
					Question: "",
					Ok:       true,
					Current:  int32(0),
					Total:    int32(0),
				}

				// Se envia y se hace las validaciones de si es que dieron finalizar el streaming de data
				// o si ocurrio un error
				err := stream.Send(questionToSend)
				// EndOfFile ---> no mas data streaming
				if err == io.EOF {
					return nil
				}
				// Un error de verdad
				if err != nil {
					return err
				}
				// Se corta el ciclo actual y se deja la app para que reciba el siguiente test
				break
			}
		}
	}
}
