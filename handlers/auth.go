package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/middleware"
	"github.com/Anwarjondev/task-management-api/models"
	"github.com/Anwarjondev/task-management-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Register a new user
// @Summary Register a new user
// @Description Register a user with username, password, and optional role
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := validate.Struct(&user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	if user.Role == "" {
		user.Role = "team_member"
	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user.Password = string(hashedpassword)
	if err = db.DB.Create(&user).Error; err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registred successfully"})
}


// Login and get JWT
// @Summary User login
// @Description Login with username and password to receive a JWT
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	err := validate.Struct(&user)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, err.Error())
	}
	var dbUser models.User
	if err := db.DB.Where("username = ?", user.Username).First(&dbUser).Error; err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dbUser.Password)); err != nil {
		utils.SendError(w, http.StatusUnauthorized, err.Error())
		return
	}
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &middleware.Claims{
		UserID: dbUser.ID,
		Role:   dbUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
