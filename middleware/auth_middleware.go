package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Anwarjondev/task-management-api/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)


var JwtKey []byte

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error with loading .env file")
	}
	JwtKey = []byte(os.Getenv("JWT_KEY")) 
}

type Claims struct {
	UserID string `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendError(w, http.StatusUnauthorized, "Unauthorized: Missing token")
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		}) 
		if err != nil || !token.Valid {
			utils.SendError(w, http.StatusUnauthorized, "Unathorized: "+err.Error())
			return
		}
		if claims.UserID == "" {
			utils.SendError(w, http.StatusUnauthorized, "Unathorized: Missing user id")
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}