package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"
	"time"
	taskapi "todoshnik/internal/api"
	"todoshnik/internal/app"
	grpcmiddleware "todoshnik/internal/grpc/middleware"
	"todoshnik/internal/grpc/pb"

	"google.golang.org/grpc"
	gogrpc "google.golang.org/grpc"

	taskgrpc "todoshnik/internal/task/grpc"
)

var logFileName string = "/api.log"

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	container := app.InitApp(logFileName)

	apiHandler := taskapi.NewAPIHandler(container)

	go func() {
		if err := apiHandler.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	grpcHandler := taskgrpc.NewHandler(container)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := gogrpc.NewServer(grpc.UnaryInterceptor(
		grpcmiddleware.AuthInterceptor(container.UserService, container.TokenService),
	))

	pb.RegisterTaskServiceServer(
		grpcServer,
		grpcHandler,
	)

	go func() {
		log.Println("gRPC server started on :50051")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := apiHandler.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}

	grpcServer.GracefulStop()

	if err := container.Cache.Close(); err != nil {
		log.Println(err)
	}

	if err := container.LogFile.Close(); err != nil {
		log.Println(err)
	}
}
