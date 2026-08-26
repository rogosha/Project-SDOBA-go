package main

import (
	"SDOBA/internal/handler"
	"SDOBA/internal/repository"
	"SDOBA/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"

	"SDOBA/internal/config"
	"SDOBA/internal/database"
	"SDOBA/internal/middleware"
	"github.com/gofiber/swagger"
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

	app.Get("/swagger/*", swagger.HandlerDefault)

	tokenService := service.NewTokenService(cfg.JWT.Secret)
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	authService := service.NewAuthService(userService, tokenService)
	userAuthHandler := handler.NewAuthHandler(authService)

	api := app.Group("/api/v1")

	users := api.Group("/users", middleware.JWTAuth(cfg.JWT.Secret))
	users.Get("/:id", userHandler.GetByID)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	auth := api.Group("/auth")
	auth.Post("/register", userAuthHandler.Register)
	auth.Post("/login", userAuthHandler.Login)

	protectedAuth := auth.Group("", middleware.JWTAuth(cfg.JWT.Secret))
	protectedAuth.Get("/me", middleware.JWTAuth(cfg.JWT.Secret), userAuthHandler.Me)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "ok",
			"database": "ok",
			"service":  cfg.App.Name,
		})
	})

	log.Fatal(app.Listen(":" + cfg.App.Port))
}
