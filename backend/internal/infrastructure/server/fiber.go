package server

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"haberbot/internal/adapter/handler"
	"haberbot/internal/infrastructure/container"
	"haberbot/internal/infrastructure/auth"
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
		AllowHeaders:     "Origin, Content-Type, Accept, X-Admin-API-Key, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Rate limiter for auth endpoints (register/login)
	authLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Çok fazla istek yapıldı. Lütfen daha sonra tekrar deneyiniz.",
			})
		},
	})

	// Register routes
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/articles", container.ArticleHandler.GetRecentArticles)
	api.Get("/articles/search", container.ArticleHandler.SearchArticles)
	api.Get("/articles/:id", container.ArticleHandler.GetArticleByID)
	api.Post("/articles/:id/summary", container.ArticleHandler.SummarizeArticle)
	api.Get("/topics", container.TopicHandler.GetTopics)
	api.Get("/topics/:slug/articles", container.TopicHandler.GetTopicArticles)

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authLimiter, container.AuthHandler.Register)
	authGroup.Post("/login", authLimiter, container.AuthHandler.Login)
	authGroup.Post("/refresh", container.AuthHandler.Refresh)
	authGroup.Post("/logout", handler.JWTMiddleware(container.Config.JWTSecret), container.AuthHandler.Logout)

	// Comments routes
	api.Get("/comments", container.CommentHandler.ListComments)
	api.Post("/comments", handler.JWTMiddleware(container.Config.JWTSecret), container.CommentHandler.CreateComment)

	// Custom Admin Gate supporting legacy API Key OR JWT role "Admin"
	adminAuth := func(c *fiber.Ctx) error {
		// 1. Try legacy API Key first
		key := c.Get("X-Admin-API-Key")
		if key != "" && key == container.Config.AdminAPIKey {
			return c.Next()
		}

		// 2. Try JWT Auth
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				claims, err := auth.ValidateToken(parts[1], container.Config.JWTSecret)
				if err == nil {
					role, ok := claims["role"].(string)
					if ok && strings.ToLower(role) == "admin" {
						c.Locals("user", claims)
						return c.Next()
					}
				}
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized admin access",
		})
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(adminAuth)
	admin.Post("/topics", container.AdminHandler.CreateTopic)
	admin.Put("/topics/:id", container.AdminHandler.UpdateTopic)
	admin.Delete("/topics/:id", container.AdminHandler.DeleteTopic)
	admin.Post("/fetch", container.AdminHandler.TriggerFetch)
	admin.Post("/process", container.AdminHandler.TriggerProcess)
	admin.Get("/articles", container.AdminHandler.ListArticlesAdmin)
	admin.Put("/articles/:id", container.AdminHandler.UpdateArticleAdmin)
	admin.Delete("/articles/:id", container.AdminHandler.DeleteArticleAdmin)

	// Health check route
	app.Get("/health", container.HealthHandler.HealthCheck)

	return app
}
