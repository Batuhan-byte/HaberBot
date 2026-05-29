package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"haberbot/internal/usecase"
)

// AuthHandler handles HTTP requests for user authentication (Register, Login, Refresh, Logout).
type AuthHandler struct {
	registerUC        *usecase.AuthRegisterUseCase
	loginUC           *usecase.AuthLoginUseCase
	refreshUC         *usecase.AuthRefreshUseCase
	logoutUC          *usecase.AuthLogoutUseCase
	jwtRefreshTTLDays int
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(
	registerUC *usecase.AuthRegisterUseCase,
	loginUC *usecase.AuthLoginUseCase,
	refreshUC *usecase.AuthRefreshUseCase,
	logoutUC *usecase.AuthLogoutUseCase,
	jwtRefreshTTLDays int,
) *AuthHandler {
	return &AuthHandler{
		registerUC:        registerUC,
		loginUC:           loginUC,
		refreshUC:         refreshUC,
		logoutUC:          logoutUC,
		jwtRefreshTTLDays: jwtRefreshTTLDays,
	}
}

// Register handles POST /api/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req usecase.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, err := h.registerUC.Execute(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "registration successful",
		"user":    user,
	})
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req usecase.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	loginResp, err := h.loginUC.Execute(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set HttpOnly Refresh Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Expires:  time.Now().AddDate(0, 0, h.jwtRefreshTTLDays),
		HTTPOnly: true,
		Secure:   false, // Set to true in production
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(loginResp)
}

// Refresh handles POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// Read HttpOnly Cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing refresh token",
		})
	}

	refreshResp, err := h.refreshUC.Execute(c.Context(), refreshToken)
	if err != nil {
		// If compromise detected or invalid, clear cookie
		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-24 * time.Hour),
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set new rotated HttpOnly Refresh Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshResp.RefreshToken,
		Expires:  time.Now().AddDate(0, 0, h.jwtRefreshTTLDays),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(refreshResp)
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Read Refresh token to clear in DB
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		// Extract claims to find user sub ID (or clear it dynamically from claims if logged in)
		// But clearing cookie is primary security action.
		// For thoroughness, we also invalidate the token in HTTP context if verified.
		if userClaims, ok := c.Locals("user").(fiber.Map); ok {
			if userID, ok := userClaims["sub"].(string); ok {
				_ = h.logoutUC.Execute(c.Context(), userID)
			}
		}
	}

	// Invalidate HttpOnly Refresh Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"message": "logged out successfully",
	})
}
