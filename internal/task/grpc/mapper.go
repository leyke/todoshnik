package grpc

import (
	"todoshnik/internal/grpc/pb"
	"todoshnik/internal/task"
)

func taskToPb(
	t *task.Task,
) *pb.TaskResponse {
	return &pb.TaskResponse{
		Id:    int64(t.ID),
		Title: t.Title,
		Done:  t.Done,
	}
}

func pbToTask(
	t *pb.TaskResponse,
) *task.Task {
	return &task.Task{
		ID:    int(t.Id),
		Title: t.Title,
		Done:  t.Done,
	}
}
