package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-unit-tests-only-32-bytes!!"

// TestValidateToken_ValidToken verifies that a freshly generated token passes validation (T4).
func TestValidateToken_ValidToken(t *testing.T) {
	token, err := GenerateAccessToken("user-123", "testuser", "User", testSecret, 15)
	require.NoError(t, err)

	claims, err := ValidateToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, "user-123", claims["sub"])
	assert.Equal(t, "testuser", claims["username"])
	assert.Equal(t, "User", claims["role"])
}

// TestValidateToken_WrongSecret verifies that tokens signed with a different secret are rejected.
// This is critical: it proves the hardcoded fallback removal (BUG-002) actually matters.
func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateAccessToken("user-123", "testuser", "User", testSecret, 15)
	require.NoError(t, err)

	_, err = ValidateToken(token, "completely-different-secret")
	assert.Error(t, err, "tokens signed with wrong secret must be rejected")
}

// TestValidateToken_ExpiredToken verifies that expired tokens are rejected (T4).
func TestValidateToken_ExpiredToken(t *testing.T) {
	// Generate a token that expired 1 minute ago
	claims := jwt.MapClaims{
		"sub":      "user-123",
		"username": "testuser",
		"role":     "User",
		"exp":      time.Now().Add(-1 * time.Minute).Unix(),
		"iat":      time.Now().Add(-2 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)

	_, err = ValidateToken(tokenStr, testSecret)
	assert.Error(t, err, "expired tokens must be rejected")
}

// TestValidateToken_MalformedToken verifies that garbage input is rejected (T4).
func TestValidateToken_MalformedToken(t *testing.T) {
	malformedTokens := []string{
		"",
		"not.a.token",
		"Bearer eyJhbGci.garbage",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.garbage.signature",
		"../../../etc/passwd",
	}

	for _, tok := range malformedTokens {
		_, err := ValidateToken(tok, testSecret)
		assert.Error(t, err, "malformed token %q should be rejected", tok)
	}
}

// TestValidateToken_WrongAlgorithm verifies algorithm confusion attacks are blocked (T4).
// An attacker could create a token with alg=none to bypass verification.
func TestValidateToken_WrongAlgorithm(t *testing.T) {
	// Craft a token with alg=none (algorithm confusion attack)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub":  "attacker",
		"role": "Admin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	// jwt.UnsafeAllowNoneSignatureType is required for signing with alg=none
	tokenStr, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	_, err = ValidateToken(tokenStr, testSecret)
	assert.Error(t, err, "tokens with alg=none must be rejected")
}

// TestGenerateAccessToken_ContainsRequiredClaims verifies all required claims are present.
func TestGenerateAccessToken_ContainsRequiredClaims(t *testing.T) {
	tokenStr, err := GenerateAccessToken("uid-999", "alice", "Admin", testSecret, 30)
	require.NoError(t, err)

	claims, err := ValidateToken(tokenStr, testSecret)
	require.NoError(t, err)

	// Verify all expected claims exist
	assert.Equal(t, "uid-999", claims["sub"])
	assert.Equal(t, "alice", claims["username"])
	assert.Equal(t, "Admin", claims["role"])
	assert.NotNil(t, claims["exp"])
	assert.NotNil(t, claims["iat"])
}

// TestGenerateRefreshToken_ValidAndRotatable verifies refresh token round-trip.
func TestGenerateRefreshToken_ValidAndRotatable(t *testing.T) {
	tokenStr, err := GenerateRefreshToken("user-456", testSecret, 7)
	require.NoError(t, err)

	claims, err := ValidateToken(tokenStr, testSecret)
	require.NoError(t, err)

	assert.Equal(t, "user-456", claims["sub"])
	// Refresh tokens should not contain role or username claims
	_, hasRole := claims["role"]
	_, hasUsername := claims["username"]
	assert.False(t, hasRole, "refresh tokens must not contain role claim")
	assert.False(t, hasUsername, "refresh tokens must not contain username claim")
}
