package main

import (
	"fmt"
	"log"
	"os"

	"chat_app/internal/database"
	"chat_app/internal/handler"
	"chat_app/internal/model"
	"chat_app/internal/repository"
	"chat_app/internal/routes"
	"chat_app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	db, err :=
		database.Connect(dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.User{},
	)

	if err != nil {
		log.Fatal(err)
	}

	userRepo :=
		repository.NewUserRepository(db)

	authService :=
		service.NewAuthService(
			userRepo,
		)

	userService :=
		service.NewUserService(
			userRepo,
		)

	authHandler :=
		handler.NewAuthHandler(
			authService,
		)

	userHandler :=
		handler.NewUserHandler(
			userService,
		)

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
		"User Service running on :8080",
	)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
