package grpc

import (
	"context"

	gogrpc "google.golang.org/grpc"

	"todoshnik/internal/grpc/pb"
	"todoshnik/internal/task"
)

type Client struct {
	client pb.TaskServiceClient
}

func NewClient(
	conn *gogrpc.ClientConn,
) *Client {
	return &Client{
		client: pb.NewTaskServiceClient(conn),
	}
}

func (c *Client) CreateTask(ctx context.Context, title string) (*task.Task, error) {
	ctx = authContext(ctx)

	response, err := c.client.CreateTask(
		ctx,
		&pb.CreateTaskRequest{
			Title: title,
		},
	)

	if err != nil {
		return nil, err
	}

	return pbToTask(response), nil
}

func (c *Client) ListTasks(ctx context.Context, status string) ([]*task.Task, error) {
	ctx = authContext(ctx)

	response, err := c.client.ListTasks(
		ctx,
		&pb.ListTasksRequest{
			Status: status,
		},
	)

	if err != nil {
		return nil, err
	}

	tasks := make([]*task.Task, 0, len(response.Tasks))

	for _, t := range response.Tasks {
		tasks = append(tasks, pbToTask(t))
	}

	return tasks, nil
}
