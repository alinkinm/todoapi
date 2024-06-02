package repository

import (
	"context"
	"fmt"
	"todoapi/internal/core"

	"strings"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	Db *sqlx.DB
}

const (
	GetAllTasks = "SELECT * FROM Task WHERE 1=1"
	CreateTask  = "INSERT INTO Task(header, descr, task_date, done) values ($1, $2, $3, $4) RETURNING id;"

	GetTaskById = "SELECT * FROM Task WHERE id= ($1);"
	UpdateTask  = "UPDATE Task SET"
	DeleteTask  = "DELETE FROM Task WHERE id = ($1);"
)

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{Db: db}
}

func (repository *TaskRepository) GetAllTasks(ctx context.Context, done string, date string, limit int, offset int) ([]*core.Task, error) {

	var queryBuilder strings.Builder
	var tasks []*core.Task
	var args []interface{}

	queryBuilder.WriteString(GetAllTasks)

	i := 1

	if done != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND done = ($%d)", i))
		args = append(args, done)
		i++
	}

	if date != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND task_date = ($%d)", i))
		args = append(args, date)
		i++
	}

	if limit != 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT ($%d)", i))
		args = append(args, limit)
		i++
	}

	if offset != 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET ($%d)", i))
		args = append(args, offset)
	}

	queryBuilder.WriteString(";")

	query := queryBuilder.String()

	rows, err := repository.Db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting all tasks: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		task := &core.Task{}
		err = rows.StructScan(task)
		if err != nil {
			return nil, fmt.Errorf("error getting all tasks: %w", err)
		}
		task.TaskDate = task.TaskDate[:10]
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (repository *TaskRepository) CreateTask(ctx context.Context, task *core.Task) (int, error) {

	id := -1

	err := repository.Db.QueryRowContext(ctx, CreateTask, task.Header, task.Description, task.TaskDate, false).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("error creating new task: %w", err)
	}

	return id, nil

}

func (repository *TaskRepository) GetTaskById(ctx context.Context, id int) (*core.Task, error) {

	task := &core.Task{}
	rows, err := repository.Db.QueryxContext(ctx, GetTaskById, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving a task by id: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(task)
		if err != nil {
			return nil, fmt.Errorf("error getting all tasks: %w", err)
		}
		task.TaskDate = task.TaskDate[:10]
	}

	return task, nil
}

func (repository *TaskRepository) UpdateTask(ctx context.Context, id int, task *core.Task) error {

	var queryBuilder strings.Builder
	var args []interface{}
	queryBuilder.WriteString(UpdateTask)

	i := 1

	if task.Header != "" {
		queryBuilder.WriteString(fmt.Sprintf(" header = $%d ", i))
		args = append(args, task.Header)
		i++
	}

	if task.Description != "" {
		queryBuilder.WriteString(fmt.Sprintf(" description = $%d ", i))
		args = append(args, task.Description)
		i++
	}

	if task.TaskDate != "" {
		queryBuilder.WriteString(fmt.Sprintf(" date = $%d", i))
		args = append(args, task.TaskDate)
		i++
	}

	if task.Done != nil {
		queryBuilder.WriteString(fmt.Sprintf(" done = $%d", i))
		args = append(args, task.Done)
		i++
	}

	queryBuilder.WriteString(fmt.Sprintf(" WHERE id = $%d ;", i))
	args = append(args, id)

	query := queryBuilder.String()

	result, err := repository.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error updating a task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating a task: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task with this id does not exist")
	}

	return nil

}

func (repository *TaskRepository) DeleteTask(ctx context.Context, id int) error {

	result, err := repository.Db.ExecContext(ctx, DeleteTask, id)

	if err != nil {
		return fmt.Errorf("error deleting task %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error updating a task: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("task with this id does not exist")
	}

	return nil
}
