package handler

import (
	"context"
	"database/sql"
	"todoapi/internal/core"

	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"time"

	"errors"

	"github.com/go-playground/validator/v10"
)

//go:generate mockery --name=TaskRepository --output=mocks --outpkg=mocks
type TaskRepository interface {
	GetAllTasks(ctx context.Context, done string, date string, limit int, offset int) ([]*core.Task, error)
	CreateTask(ctx context.Context, task *core.Task) (int, error)
	GetTaskById(ctx context.Context, id int) (*core.Task, error)
	UpdateTask(ctx context.Context, id int, task *core.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type TaskHandler struct {
	taskRepository TaskRepository
}

func NewTaskHandler(repository TaskRepository) *TaskHandler {
	return &TaskHandler{taskRepository: repository}
}

func (handler *TaskHandler) InitRoutes(app *fiber.App) {

	app.Get("/tasks", handler.GetAllTasks)
	app.Post("/tasks", handler.CreateTask)
	app.Get("/tasks/:id", handler.GetTaskById)
	app.Patch("/tasks/:id", handler.UpdateTask)
	app.Delete("tasks/:id", handler.DeleteTask)
}

// Get All Tasks
// @Summary Get all created tasks
// @Tags tasks
// @Description Returns list of tasks
// @Produce json
// @Param done query boolean false "status filter"
// @Param date query string false "date filter"
// @Param pageSize query int false "page capacity (in tasks)"
// @Param page query int false "page number"
// @Success 200 {object} map[string][]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (handler *TaskHandler) GetAllTasks(ctx *fiber.Ctx) error {

	done := ctx.Query("done")
	if done != "" && done != "false" && done != "true" {
		log.Info("error while parsing status param")
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error ": "Incorrect Task Status Parameter",
			})
	}

	date := ctx.Query("date")
	if date != "" {
		_, err := isValidDate(date)
		if err != nil {
			log.Info("error while parsing date param")
			return ctx.Status(http.StatusBadRequest).JSON(
				fiber.Map{
					"Error ": "Incorrect Date Format",
				})
		}
	}
	pageSize := ctx.Query("pageSize")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil && pageSize != "" {
		log.Info("error while parsing page size param")
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error ": "Incorrect input for page size parameter",
			})
	}

	page := ctx.Query("page")
	if page == "" && pageSize != "" {
		page = "1"
	}
	if pageSize == "" && page != "" {
		pageSizeInt = 5
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil && page != "" {
		log.Info("error while parsing page param")
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error ": "Incorrect input for page parameter",
			})
	}
	log.Info("params parsed")

	tasks, err := handler.taskRepository.GetAllTasks(ctx.UserContext(), done, date, pageSizeInt, (pageInt-1)*pageSizeInt)
	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"Error": "Internal Server Error",
			})
	}

	log.Info("got tasks")

	if tasks == nil {
		tasks = []*core.Task{}
	}

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"Tasks": tasks,
		})

}

// Create new task
// @Summary Create new task
// @Tags tasks
// @Description Create new task
// @Produce json
// @Param header body string true "header"
// @Param descr body string true "description"
// @Param date body string true "date YYYY-MM-DD"
// @Success 201 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (handler *TaskHandler) CreateTask(ctx *fiber.Ctx) error {

	task := &core.Task{}

	if err := ctx.BodyParser(task); err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error parsing json": "Incorrect input",
			})
	}

	date := task.TaskDate
	if date != "" {
		_, err := isValidDate(date)
		if err != nil {
			log.Info(err.Error())
			return ctx.Status(http.StatusBadRequest).JSON(
				fiber.Map{
					"Error ": "Incorrect date format",
				})
		}
	}

	log.Info("Input json parsed")

	validate := validator.New()
	if err := validate.Struct(task); err != nil {
		log.Info(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Validation error": "Invalid input",
		})
	}

	log.Info("Input validated")

	id, err := handler.taskRepository.CreateTask(ctx.UserContext(), task)

	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": "Internal Server Error",
			})
	}

	log.Info("task saved")

	return ctx.Status(http.StatusCreated).JSON(
		fiber.Map{
			"Task id": id,
		})
}

// Get Task
// @Summary Get task by id
// @Tags tasks
// @Description Returns task if it exists
// @Produce json
// @Param id path int true "task id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func (handler *TaskHandler) GetTaskById(ctx *fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Bad Request",
			})
	}

	log.Info("task id parsed")

	task, err := handler.taskRepository.GetTaskById(ctx.UserContext(), id)
	if err != nil {
		log.Info(err.Error())
		e := "Internal Server Error"
		if errors.Unwrap(err) == sql.ErrNoRows {
			e = "No tasks with this id"
		}
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": e,
			})
	}

	if (core.Task{}) == *task {
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error": "No tasks with this id",
			})
	}

	log.Info("got a task")

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"Task": task,
		})

}

// Update Task
// @Summary Update task by id
// @Tags tasks
// @Description Update 1+ properties of an existing task
// @Produce json
// @Param id path int true "task id"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [patch]
func (handler *TaskHandler) UpdateTask(ctx *fiber.Ctx) error {

	task := &core.Task{}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error": "Bad Request",
			})
	}

	log.Info("task id parsed")

	if err := ctx.BodyParser(task); err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error parsing json": "Invalid input",
			})
	}

	_, err = isUpdateValid(task)
	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"Error ": "Invalid task format",
			})
	}

	date := task.TaskDate
	if date != "" {
		_, err := isValidDate(date)
		if err != nil {
			log.Info(err.Error())
			return ctx.Status(http.StatusBadRequest).JSON(
				fiber.Map{
					"Error ": "Incorrect date format",
				})
		}
	}

	log.Info("Input json parsed")

	err = handler.taskRepository.UpdateTask(ctx.UserContext(), id, task)
	if err != nil {
		log.Info(err.Error())
		e := "Internal Server Error"
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": e,
			})
	}

	log.Info("updated a task")

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"Updated task with id": id,
		})

}

// Delete Task
// @Summary Delete task by id
// @Tags tasks
// @Description Delete task
// @Produce json
// @Param id path int true "task id"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (handler *TaskHandler) DeleteTask(ctx *fiber.Ctx) error {

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		log.Info(err.Error())
		return ctx.Status(http.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Bad Request",
			})
	}

	log.Info("task id parsed")

	err = handler.taskRepository.DeleteTask(ctx.UserContext(), id)
	if err != nil {
		log.Info(err.Error())
		e := "Internal Server Error"

		if err.Error() == "task with this id does not exist" {
			e = err.Error()
		}
		return ctx.Status(http.StatusInternalServerError).JSON(
			fiber.Map{
				"error": e,
			})
	}

	log.Info("got a task")

	return ctx.Status(http.StatusOK).JSON(
		fiber.Map{
			"Deleted task with id": id,
		})
}

func isValidDate(dateStr string) (bool, error) {

	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	log.Info(t)
	if err != nil {
		return false, err
	}
	formattedDate := t.Format(layout)
	log.Info(formattedDate)
	if len(dateStr) == len(formattedDate) {
		return true, nil
	}

	return false, err
}

func isUpdateValid(task *core.Task) (bool, error) {

	if (core.Task{}) == *task {
		return false, errors.New("empty task received")
	}

	if task.Done == nil {
		return false, errors.New("empty task received")
	}

	return true, nil

}
