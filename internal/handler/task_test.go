// package handler

// import (
// 	"testing"

// 	"todoapi/internal/core"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/mock"

// 	"context"
// )

// type MockTaskRepository struct {
// 	mock.Mock
// }

// func (m *MockTaskRepository) GetAllTasks(ctx context.Context, done string, date string, limit int, offset int) ([]*core.Task, error) {
// 	args := m.Called(ctx, done, date, limit, offset)
// 	return args.Get(0).([]*core.Task), args.Error(1)
// }

// func TestTaskHandler_GetAllTasks(t *testing.T) {

// 	mockRepo := new(MockTaskRepository)
// 	mockRepo.On("GetAll").Return([]*core.Task{{ID: 1, Name: "Test Task"}}, nil)
// 	taskHandler := NewTaskHandler(mockRepo)

// 	type args struct {
// 		ctx *fiber.Ctx
// 	}
// 	tests := []struct {
// 		name    string
// 		handler *TaskHandler
// 		args    args
// 		wantErr bool
// 	}{

// 		{
// 			name: "all tasks success",
// 			handler: &TaskHandler{
// 				taskRepository: mockRepo,
// 			},
// 		},

// 		{
// 			name: "Успешное получение всех задач",
// 			handler: &TaskHandler{
// 				Service: &MockTaskService{},
// 			},
// 			args: args{
// 				ctx: app.AcquireCtx(&fiber.Ctx{}),
// 			},
// 			setupMock: func(s *MockTaskService) {
// 				s.On("GetAll").Return([]Task{{ID: 1, Name: "Test Task"}}, nil)
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := tt.handler.GetAllTasks(tt.args.ctx); (err != nil) != tt.wantErr {
// 				t.Errorf("TaskHandler.GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

package handler_test
