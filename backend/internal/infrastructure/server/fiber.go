package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"haberbot/internal/infrastructure/container"
)

// SetupFiberServer creates and configures a new Fiber app, registering all routes.
func SetupFiberServer(container *container.Container) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "HaberBot API",
	})

	// Add standard middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Configure CORS for the frontend origin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     container.Config.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, X-Admin-API-Key",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Register routes
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/articles", container.ArticleHandler.GetRecentArticles)
	api.Get("/articles/search", container.ArticleHandler.SearchArticles)
	api.Get("/articles/:id", container.ArticleHandler.GetArticleByID)
	api.Post("/articles/:id/summary", container.ArticleHandler.SummarizeArticle)
	api.Get("/topics", container.TopicHandler.GetTopics)
	api.Get("/topics/:slug/articles", container.TopicHandler.GetTopicArticles)

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(container.AdminHandler.AuthMiddleware())
	admin.Post("/topics", container.AdminHandler.CreateTopic)
	admin.Put("/topics/:id", container.AdminHandler.UpdateTopic)
	admin.Delete("/topics/:id", container.AdminHandler.DeleteTopic)
	admin.Post("/fetch", container.AdminHandler.TriggerFetch)
	admin.Post("/process", container.AdminHandler.TriggerProcess)

	// Health check route
	app.Get("/health", container.HealthHandler.HealthCheck)

	return app
}
