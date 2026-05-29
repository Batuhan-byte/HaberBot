package server

import (
	"crypto/subtle"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"haberbot/internal/adapter/handler"
	"haberbot/internal/infrastructure/auth"
	"haberbot/internal/infrastructure/container"
)

// SetupFiberServer creates and configures a new Fiber app, registering all routes.
func SetupFiberServer(container *container.Container) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "HaberBot API",
		// BUG-017: Do not expose stack traces or error details to clients.
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": "an internal error occurred"})
		},
	})

	// Panic recovery
	app.Use(recover.New())
	app.Use(logger.New())

	// BUG-011: Security headers middleware (OWASP recommended)
	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		if container.Config.IsProduction {
			c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		}
		return c.Next()
	})

	// Configure CORS for the frontend origin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     container.Config.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, X-Admin-API-Key, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Rate limiter for auth endpoints (register/login) — strict
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

	// BUG-009: Refresh endpoint also needs its own rate limiter
	refreshLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Çok fazla yenileme isteği yapıldı.",
			})
		},
	})

	// Static file serving for uploads (avatars)
	app.Static("/avatars", "./public/avatars")
	app.Static("/public/avatars", "./public/avatars")

	// Register routes
	api := app.Group("/api/v1")

	// Public routes
	api.Get("/articles", container.ArticleHandler.GetRecentArticles)
	api.Get("/articles/search", container.ArticleHandler.SearchArticles)
	api.Get("/articles/:id", container.ArticleHandler.GetArticleByID)

	// BUG-008 & Phase 9: Summary endpoint. Made public by explicit user request to attract guest readers.
	// We support both GET (requested by Phase 9 spec) and POST (backward compatibility).
	api.Get("/articles/:id/summary",
		container.ArticleHandler.SummarizeArticle,
	)
	api.Post("/articles/:id/summary",
		container.ArticleHandler.SummarizeArticle,
	)

	api.Get("/topics", container.TopicHandler.GetTopics)
	api.Get("/topics/:slug/articles", container.TopicHandler.GetTopicArticles)

	// User Profiles (Public Profile View) - optional JWT to detect favorited status
	api.Get("/users/:username", func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				claims, err := auth.ValidateToken(parts[1], container.Config.JWTSecret)
				if err == nil {
					c.Locals("user", claims)
				}
			}
		}
		return container.ProfileHandler.GetProfile(c)
	})

	// Auth routes
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authLimiter, container.AuthHandler.Register)
	authGroup.Post("/login", authLimiter, container.AuthHandler.Login)
	// BUG-009: refresh now rate-limited
	authGroup.Post("/refresh", refreshLimiter, container.AuthHandler.Refresh)
	authGroup.Post("/logout", handler.JWTMiddleware(container.Config.JWTSecret), container.AuthHandler.Logout)

	// Comments routes
	api.Get("/comments", container.CommentHandler.ListComments)
	api.Post("/comments", handler.JWTMiddleware(container.Config.JWTSecret), container.CommentHandler.CreateComment)

	// Profile authenticated routes
	profileGroup := api.Group("/users")
	profileGroup.Use(handler.JWTMiddleware(container.Config.JWTSecret))
	profileGroup.Put("/profile", container.ProfileHandler.UpdateProfile)
	profileGroup.Post("/:userId/favorites", container.ProfileHandler.AddFavorite)
	profileGroup.Delete("/:userId/favorites", container.ProfileHandler.RemoveFavorite)
	profileGroup.Get("/me/favorites", container.ProfileHandler.ListFavorites)

	// Profile reports
	api.Post("/profile-reports", handler.JWTMiddleware(container.Config.JWTSecret), container.ProfileHandler.ReportProfile)

	// BUG-010: Admin API key comparison uses constant-time compare to prevent timing attacks.
	adminAuth := func(c *fiber.Ctx) error {
		// 1. Try legacy API Key first (constant-time comparison)
		key := c.Get("X-Admin-API-Key")
		configKey := container.Config.AdminAPIKey
		if key != "" && configKey != "" &&
			subtle.ConstantTimeCompare([]byte(key), []byte(configKey)) == 1 {
			return c.Next()
		}

		// 2. Try JWT Auth with Admin role
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
