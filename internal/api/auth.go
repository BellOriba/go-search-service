package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/BellOriba/go-search-service/internal/products"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(repo *products.PostgresRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req products.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		user, err := repo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID.String(),
			"role":    user.Role,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products.LoginResponse{Token: tokenString})
	}
}

