package repository

import (
	"context"

	"github.com/JuanDiegoCastellanos/grpc/models"
)

type Repository interface {
	GetStudent(ctx context.Context, id string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error
	GetTest(ctx context.Context, id string) (*models.Test, error)
	SetTest(ctx context.Context, test *models.Test) error
	SetQuestion(ctx context.Context, question *models.Question) error
	SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error)
	GetQuestionPerTest(ctx context.Context, testId string) ([]*models.Question, error)
	SetAnswer(ctx context.Context, answer *models.Answer) error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func SetStudent(ctx context.Context, student *models.Student) error {
	return implementation.SetStudent(ctx, student)
}
func GetStudent(ctx context.Context, id string) (*models.Student, error) {
	return implementation.GetStudent(ctx, id)
}
func GetTest(ctx context.Context, id string) (*models.Test, error) {
	return implementation.GetTest(ctx, id)
}
func SetTest(ctx context.Context, test *models.Test) error {
	return implementation.SetTest(ctx, test)
}
func SetQuestion(ctx context.Context, question *models.Question) error {
	return implementation.SetQuestion(ctx, question)
}
func SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return implementation.SetEnrollment(ctx, enrollment)
}
func GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {
	return implementation.GetStudentsPerTest(ctx, testId)
}
func GetQuestionPerTest(ctx context.Context, testId string) ([]*models.Question, error) {
	return implementation.GetQuestionPerTest(ctx, testId)
}
func SetAnswer(ctx context.Context, answer *models.Answer) error {
	return implementation.SetAnswer(ctx, answer)
}
