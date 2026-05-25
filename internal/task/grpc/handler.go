package grpc

import (
	"context"
	"todoshnik/internal/app"
	"todoshnik/internal/grpc/pb"
	"todoshnik/internal/task"

	authcontext "todoshnik/internal/auth/context"
	apperrors "todoshnik/internal/errors"
)

type Handler struct {
	pb.UnimplementedTaskServiceServer
	service *task.Service
}

func NewHandler(container *app.App) *Handler {
	return &Handler{
		service: container.TaskService,
	}
}

func (h *Handler) CreateTask(
	ctx context.Context,
	req *pb.CreateTaskRequest,
) (*pb.TaskResponse, error) {
	userID, ok := authcontext.GetUserID(ctx)
	if !ok {
		return nil, apperrors.ErrUnAuth
	}

	newTask, err := h.service.Add(
		ctx,
		req.Title,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return taskToPb(newTask), nil
}

func (h *Handler) ListTasks(
	ctx context.Context,
	req *pb.ListTasksRequest,
) (*pb.ListTasksResponse, error) {
	userID, ok := authcontext.GetUserID(ctx)
	if !ok {
		return nil, apperrors.ErrUnAuth
	}

	scope := getScope(userID)

	tasks, err := h.service.List(
		ctx,
		task.TaskFilter{
			Scope: scope,
		})

	if err != nil {
		return nil, err
	}

	response := &pb.ListTasksResponse{
		Tasks: make([]*pb.TaskResponse, 0, len(tasks)),
	}

	for _, t := range tasks {
		response.Tasks = append(
			response.Tasks,
			taskToPb(t),
		)
	}

	return response, nil
}
