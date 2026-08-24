package main

import (
	"SDOBA/internal/handler"
	"SDOBA/internal/repository"
	"SDOBA/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"

	"SDOBA/internal/config"
	"SDOBA/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := service.NewAuthService(userService)
	userHandler := handler.NewUserHandler(userService)
	userAuthHandler := handler.NewAuthHandler(authService)

	api := app.Group("/api/v1")

	users := api.Group("/users")
	auth := api.Group("/auth")

	auth.Post("/", userAuthHandler.Register)
	auth.Post("/login", userAuthHandler.Login)
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "ok",
			"service":  cfg.App.Name,
		})
	})

	log.Fatal(app.Listen(":" + cfg.App.Port))
}
