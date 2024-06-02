package handler_test

// import (
// 	"testing"

// 	"todoapi/internal/core"
// 	"todoapi/internal/handler"
// 	"todoapi/internal/handler/mocks"

// 	"github.com/gofiber/fiber/v2"
// )

// func TestTaskHandler_GetAllTasks(t *testing.T) {

// 	mockRepo := new(mocks.MockTaskRepository)

// 	mockRepo.On("GetAll").Return([]*core.Task{
// 		{},
// 	}, nil)

// 	taskHandler := handler.NewTaskHandler(mockRepo)

// 	type args struct {
// 		ctx *fiber.Ctx
// 	}
// 	tests := []struct {
// 		name    string
// 		handler *handler.TaskHandler
// 		args    args
// 		wantErr bool
// 	}{

// 		{
// 			name:    "all tasks success",
// 			handler: taskHandler,
// 			args: args {
// 				done:  "",
// 				date:  "",
// 				pageSize: "",
// 				page:  "",
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
