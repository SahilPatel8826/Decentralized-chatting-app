package main

import (
	"fmt"
	"log"
	"os"

	pb "chat_service/pb"

	"chat_service/internal/client"
	"chat_service/internal/database"
	"chat_service/internal/handler"
	"chat_service/internal/model"
	"chat_service/internal/repository"
	"chat_service/internal/routes"
	"chat_service/internal/service"
	"chat_service/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
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
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Message{},
	)

	if err != nil {
		log.Fatal(err)
	}

	// ==========================
	// WebSocket Hub
	// ==========================

	hub := websocket.NewHub()

	go hub.Run()

	// ==========================
	// Repositories
	// ==========================

	conversationRepo :=
		repository.NewConversationRepository(
			db,
		)

	messageRepo :=
		repository.NewMessageRepository(
			db,
		)

	// ==========================
	// gRPC User Service Client
	// ==========================

	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	grpcClient :=
		pb.NewUserServiceClient(
			conn,
		)

	userClient :=
		client.NewUserClient(
			grpcClient,
		)

	// ==========================
	// Services
	// ==========================

	conversationService :=
		service.NewConversationService(
			conversationRepo,
			userClient,
		)

	messageService :=
		service.NewMessageService(
			messageRepo,
			conversationRepo,
			hub,
		)

	// ==========================
	// Handlers
	// ==========================

	conversationHandler :=
		handler.NewConversationHandler(
			conversationService,
		)

	messageHandler :=
		handler.NewMessageHandler(
			messageService,
		)

	wsHandler :=
		websocket.NewHandler(
			hub,
		)

	// ==========================
	// Router
	// ==========================

	router := gin.Default()

	routes.ConversationRoutes(
		router,
		conversationHandler,
	)

	routes.MessageRoutes(
		router,
		messageHandler,
	)

	routes.WebSocketRoutes(
		router,
		wsHandler,
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
		"Chat Service running on :8081",
	)

	if err := router.Run(
		":8081",
	); err != nil {

		log.Fatal(err)
	}

}
