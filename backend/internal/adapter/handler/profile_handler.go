package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"haberbot/internal/usecase"
)

// ProfileHandler implements HTTP handlers for the user profile domain.
type ProfileHandler struct {
	getProfileUC     *usecase.GetProfileUseCase
	updateProfileUC  *usecase.UpdateProfileUseCase
	addFavoriteUC    *usecase.AddFavoriteUseCase
	removeFavoriteUC *usecase.RemoveFavoriteUseCase
	listFavoritesUC  *usecase.ListFavoritesUseCase
	reportProfileUC  *usecase.ReportProfileUseCase
}

// NewProfileHandler creates a new ProfileHandler instance.
func NewProfileHandler(
	getProfileUC *usecase.GetProfileUseCase,
	updateProfileUC *usecase.UpdateProfileUseCase,
	addFavoriteUC *usecase.AddFavoriteUseCase,
	removeFavoriteUC *usecase.RemoveFavoriteUseCase,
	listFavoritesUC *usecase.ListFavoritesUseCase,
	reportProfileUC *usecase.ReportProfileUseCase,
) *ProfileHandler {
	return &ProfileHandler{
		getProfileUC:     getProfileUC,
		updateProfileUC:  updateProfileUC,
		addFavoriteUC:    addFavoriteUC,
		removeFavoriteUC: removeFavoriteUC,
		listFavoritesUC:  listFavoritesUC,
		reportProfileUC:  reportProfileUC,
	}
}

// GetProfile handles GET /api/v1/users/:username
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username is required"})
	}

	// Try extracting userID if logged in
	currentUserID := ""
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			// Extract JWT secret from c.App().Config() or just let the caller decode it
			// We can decode without validation or with validation. Since jwt claims might be present:
			if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
				if sub, ok := claims["sub"].(string); ok {
					currentUserID = sub
				}
			}
		}
	}

	profile, err := h.getProfileUC.Execute(c.Context(), username, currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "kullanıcı bulunamadı"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(profile)
}

// UpdateProfile handles PUT /api/v1/users/profile (multipart/form-data)
func (h *ProfileHandler) UpdateProfile(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized: missing identifier"})
	}

	bio := c.FormValue("bio")
	presetAvatar := c.FormValue("preset_avatar") // If client chose a preset avatar

	var avatarURL *string
	if presetAvatar != "" {
		// Verify valid preset URL structure
		urlPath := fmt.Sprintf("/avatars/presets/%s.png", filepath.Base(presetAvatar))
		avatarURL = &urlPath
	}

	// Handle custom file upload if present
	fileHeader, err := c.FormFile("avatar")
	if err == nil && fileHeader != nil {
		// Security Check 1: File size limit (2MB)
		if fileHeader.Size > 2*1024*1024 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Avatar boyutu en fazla 2 MB olabilir."})
		}

		// Security Check 2: Allowed extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Yalnızca JPG, PNG ve WebP dosyaları kabul edilir."})
		}

		// Open file to verify content type (MIME type defense-in-depth)
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Dosya okunamadı"})
		}
		defer file.Close()

		// Read first 512 bytes for type detection
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Dosya analizi başarısız"})
		}

		contentType := http.DetectContentType(buffer[:n])
		if !strings.HasPrefix(contentType, "image/") || 
			(contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz resim dosyası biçimi."})
		}

		// Save the file to public/avatars/ folder
		uploadDir := "./public/avatars"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Avatar dizini oluşturulamadı"})
		}

		// Create unique filename user-{userID}.ext
		filename := fmt.Sprintf("user-%s%s", userID, ext)
		destPath := filepath.Join(uploadDir, filename)

		// Delete existing files for this user with other extensions to cleanup
		extensions := []string{".jpg", ".jpeg", ".png", ".webp"}
		for _, e := range extensions {
			oldPath := filepath.Join(uploadDir, fmt.Sprintf("user-%s%s", userID, e))
			_ = os.Remove(oldPath) // ignore error if doesn't exist
		}

		// Reset seek and save
		_, _ = file.Seek(0, 0)
		out, err := os.Create(destPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Dosya kaydedilemedi"})
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Dosya yazma hatası"})
		}

		// URL paths are served under /avatars/
		urlPath := fmt.Sprintf("/avatars/%s", filename)
		avatarURL = &urlPath
	}

	req := usecase.UpdateProfileRequest{
		UserID:    userID,
		Bio:       &bio,
		AvatarURL: avatarURL,
	}

	updatedUser, err := h.updateProfileUC.Execute(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedUser)
}

// AddFavorite handles POST /api/v1/users/:userId/favorites
func (h *ProfileHandler) AddFavorite(c *fiber.Ctx) error {
	favoriteUserID := c.Params("userId")
	if favoriteUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "userId is required"})
	}

	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := h.addFavoriteUC.Execute(c.Context(), userID, favoriteUserID)
	if err != nil {
		if strings.Contains(err.Error(), "zaten") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":          "Kullanıcı favorilerinize eklendi.",
		"favorite_user_id": favoriteUserID,
	})
}

// RemoveFavorite handles DELETE /api/v1/users/:userId/favorites
func (h *ProfileHandler) RemoveFavorite(c *fiber.Ctx) error {
	favoriteUserID := c.Params("userId")
	if favoriteUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "userId is required"})
	}

	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err := h.removeFavoriteUC.Execute(c.Context(), userID, favoriteUserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Favori kaydı bulunamadı."})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListFavorites handles GET /api/v1/users/me/favorites
func (h *ProfileHandler) ListFavorites(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	favorites, total, err := h.listFavoritesUC.Execute(c.Context(), userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	c.Set("X-Total-Count", strconv.Itoa(total))
	return c.JSON(fiber.Map{
		"favorites": favorites,
		"total":     total,
	})
}

// ReportProfile handles POST /api/v1/profile-reports
func (h *ProfileHandler) ReportProfile(c *fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	type ReportPayload struct {
		ReportedUserID string `json:"reported_user_id"`
		Reason         string `json:"reason"`
		Comment        string `json:"comment"`
	}

	var payload ReportPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz şikayet verisi"})
	}

	req := usecase.ReportProfileRequest{
		ReporterUserID: userID,
		ReportedUserID: payload.ReportedUserID,
		Reason:         payload.Reason,
		Comment:        payload.Comment,
	}

	err := h.reportProfileUC.Execute(c.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "zaten") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Şikayetiniz başarıyla iletildi. Moderasyon ekibi inceleyecektir.",
	})
}
