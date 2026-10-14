package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"chat_app/internal/database"
	grpcserver "chat_app/internal/grpc"
	"chat_app/internal/handler"
	"chat_app/internal/model"
	"chat_app/internal/repository"
	"chat_app/internal/routes"
	"chat_app/internal/service"
	pb "chat_app/pb"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {

	// .env is optional: real environment variables win.
	for _, path := range []string{".env", "../../.env"} {
		_ = godotenv.Load(path)
	}

	if os.Getenv("JWT_ACCESS_SECRET") == "" ||
		os.Getenv("JWT_REFRESH_SECRET") == "" {
		log.Fatal("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be set")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := database.Connect(dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// Repositories
	// ==========================

	userRepo :=
		repository.NewUserRepository(
			db,
		)

	// ==========================
	// Services
	// ==========================

	authService :=
		service.NewAuthService(
			userRepo,
		)

	userService :=
		service.NewUserService(
			userRepo,
		)

	// ==========================
	// HTTP Handlers
	// ==========================

	authHandler :=
		handler.NewAuthHandler(
			authService,
		)

	userHandler :=
		handler.NewUserHandler(
			userService,
		)

	// ==========================
	// gRPC Server
	// ==========================

	userGrpcServer :=
		grpcserver.NewUserServer(
			userRepo,
		)

	lis, err :=
		net.Listen(
			"tcp",
			":50051",
		)

	if err != nil {
		log.Fatal(err)
	}

	grpcServer :=
		grpc.NewServer()

	pb.RegisterUserServiceServer(
		grpcServer,
		userGrpcServer,
	)

	go func() {

		log.Println(
			"gRPC Server running on :50051",
		)

		if err :=
			grpcServer.Serve(
				lis,
			); err != nil {

			log.Fatal(err)
		}
	}()

	// ==========================
	// HTTP Server
	// ==========================

	router := gin.Default()

	routes.AuthRoutes(
		router,
		authHandler,
	)

	routes.UserRoutes(
		router,
		userHandler,
	)

	router.GET(
		"/health",
		func(c *gin.Context) {

			c.JSON(
				200,
				gin.H{
					"status": "ok",
				},
			)
		},
	)

	log.Println(
		"HTTP Server running on :8080",
	)

	if err :=
		router.Run(
			":8080",
		); err != nil {

		log.Fatal(err)
	}

}
