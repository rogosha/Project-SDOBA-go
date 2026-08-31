package main

import (
	"SDOBA/internal/handler"
	"SDOBA/internal/repository"
	"SDOBA/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"

	"SDOBA/docs"

	"SDOBA/internal/config"
	"SDOBA/internal/database"
	"SDOBA/internal/middleware"

	"github.com/gofiber/swagger"
)

func main() {
	_ = docs.SwaggerInfo

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	app := fiber.New()

	app.Get("/swagger/*", swagger.HandlerDefault)

	tokenService := service.NewTokenService(cfg.JWT.Secret)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userService, tokenService)
	userAuthHandler := handler.NewAuthHandler(authService)

	conversationRepository := repository.NewConversationRepository(db)
	conversationService := service.NewConversationService(conversationRepository)
	conversationHandler := handler.NewConversationHandler(conversationService)

	messageRepository := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepository, conversationRepository)
	messageHandler := handler.NewMessageHandler(messageService)

	api := app.Group("/api/v1")

	users := api.Group("/users", middleware.JWTAuth(cfg.JWT.Secret))
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	auth := api.Group("/auth")
	auth.Post("/register", userAuthHandler.Register)
	auth.Post("/login", userAuthHandler.Login)

	protectedAuth := auth.Group("", middleware.JWTAuth(cfg.JWT.Secret))
	protectedAuth.Get("/me", userAuthHandler.Me)

	conversations := api.Group("/conversations", middleware.JWTAuth(cfg.JWT.Secret))

	conversations.Post("/", conversationHandler.Create)
	conversations.Get("/", conversationHandler.GetUserConversations)
	conversations.Get("/:id", conversationHandler.GetByID)
	conversations.Post("/:id/messages", messageHandler.Create)
	conversations.Get("/:id/messages", messageHandler.GetByConversationID)

	// Health
	// @Summary Проверка состояния сервиса
	// @Description Проверяет доступность API и базы данных
	// @Tags Health
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		dbStatus := "ok"
		redisStatus := "ok"
		status := fiber.StatusOK

		if err := sqlDB.Ping(); err != nil {
			dbStatus = "error"
			status = fiber.StatusServiceUnavailable
		}

		if err := redisClient.Ping(c.UserContext()).Err(); err != nil {
			redisStatus = "error"
			status = fiber.StatusServiceUnavailable
		}

		serviceStatus := "ok"
		if status != fiber.StatusOK {
			serviceStatus = "error"
		}

		return c.Status(status).JSON(fiber.Map{
			"status":   serviceStatus,
			"database": dbStatus,
			"redis":    redisStatus,
			"service":  cfg.App.Name,
		})
	})

	log.Fatal(app.Listen(":" + cfg.App.Port))
}
