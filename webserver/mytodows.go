package webserver

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Task model and storage
type (
	Task struct {
		ID          int64
		Description string
		Deadline    int64
	}

	// Storage
	TaskStorageInMemory struct {
		tasks map[int64]Task
	}
)

var taskIdCounter int64 = 1

func (s *TaskStorageInMemory) Create(t Task) (int64, error) {
	t.ID = taskIdCounter
	taskIdCounter++

	s.tasks[t.ID] = t

	return t.ID, nil
}

func (s *TaskStorageInMemory) List() ([]Task, error) {
	tasks := make([]Task, 0, len(s.tasks))

	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (s *TaskStorageInMemory) Read(id int64) (Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, errors.New("Task with provided ID not found")
	}

	return task, nil
}

func (s *TaskStorageInMemory) Update(id int64, upd PatchTaskRequest) (Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, errors.New("Task with provided ID not found")
	}

	if upd.Description != "" {
		task.Description = upd.Description
	}
	if upd.Deadline != 0 {
		task.Deadline = upd.Deadline
	}

	s.tasks[task.ID] = task

	return task, nil
}

func (s *TaskStorageInMemory) Delete(id int64) error {
	task, ok := s.tasks[id]
	if !ok {
		return errors.New("Task with provided ID not found")
	}

	delete(s.tasks, task.ID)

	return nil
}

// Task Creation
type (
	CreateTaskRequest struct {
		Description string `json:"description"`
		Deadline    int64  `json:"deadline"`
	}

	CreateTaskResponse struct {
		ID int64 `json:"id"`
	}
)

// Tasks Reading
type (
	ListTasksResponse struct {
		Tasks []Task `json:"tasks"`
	}

	GetTaskRequest struct {
		ID int64 `json:"id"`
	}

	GetTaskResponse struct {
		Task
	}
)

// Task Updating
type (
	PatchTaskRequest struct {
		Description string `json:"description"`
		Deadline    int64  `json:"deadline"`
	}

	PatchTaskResponse struct {
		Task
	}
)

func StartToDoServer() {
	webApp := fiber.New()
	storage := &TaskStorageInMemory{
		tasks: make(map[int64]Task),
	}

	// Create new task
	webApp.Post("/tasks", func(ctx *fiber.Ctx) error {
		var req CreateTaskRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		id, err := storage.Create(Task{
			Description: req.Description,
			Deadline:    req.Deadline,
		})
		if err != nil {
			return fmt.Errorf("creation in storage: %w", err)
		}

		return ctx.JSON(CreateTaskResponse{ID: id})
	})

	// Get list of all tasks
	webApp.Get("/tasks", func(ctx *fiber.Ctx) error {
		tasks, err := storage.List()
		if err != nil {
			return fmt.Errorf("list all tasks from storage: %w", err)
		}

		return ctx.JSON(ListTasksResponse{Tasks: tasks})
	})

	const taskIdUnknown = "unknown"
	// Get task with id
	webApp.Get("/tasks/:id", func(ctx *fiber.Ctx) error {
		taskIdParam := ctx.Params("id", taskIdUnknown)
		if taskIdParam == taskIdUnknown {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
		if err != nil {
			return fmt.Errorf("convert ID to string: %w", err)
		}

		task, err := storage.Read(taskId)
		if err != nil {
			return fmt.Errorf("read task with provided id: %w", err)
		}

		return ctx.JSON(GetTaskResponse{Task: task})
	})

	webApp.Patch("/tasks/:id", func(ctx *fiber.Ctx) error {
		taskIdParam := ctx.Params("id", taskIdUnknown)
		if taskIdParam == taskIdUnknown {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
		if err != nil {
			return fmt.Errorf("convert ID to string: %w", err)
		}

		var req PatchTaskRequest
		if err := ctx.BodyParser(&req); err != nil {
			return fmt.Errorf("body parser: %w", err)
		}

		updatedTask, err := storage.Update(taskId, req)
		if err != nil {
			return fmt.Errorf("patch task with provided id: %w", err)
		}

		return ctx.JSON(PatchTaskResponse{updatedTask})
	})

	webApp.Delete("/tasks/:id", func(ctx *fiber.Ctx) error {
		taskIdParam := ctx.Params("id", taskIdUnknown)
		if taskIdParam == taskIdUnknown {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		taskId, err := strconv.ParseInt(taskIdParam, 10, 64)
		if err != nil {
			return fmt.Errorf("convert ID to string: %w", err)
		}

		if err := storage.Delete(taskId); err != nil {
			return fmt.Errorf("delete task with provided id: %w", err)
		}

		return ctx.SendStatus(fiber.StatusOK)
	})

	port := "8080"
	logrus.Fatal(webApp.Listen(":" + port))
}
