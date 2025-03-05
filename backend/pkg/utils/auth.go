package utils

import (
	"net/http"
	"notes/internal/configs"
	"notes/pkg/date"
	"notes/pkg/validations"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (*string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, validations.ErrHashingPwd
	}
	hashed := string(hash)
	return &hashed, nil
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateJWT(userID uint) (string, error) {
	jwtExpirationDuration := time.Minute * time.Duration(configs.GetInt("JWT_TIME_COUNT", 60))
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     date.ArgentinaTimeNow().Add(jwtExpirationDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.GetString("JWT_STRING", "")))
}

// DEBUG START

func SetJWTAsCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "notes_jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   309600,
	})
}
