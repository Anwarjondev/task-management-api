package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers lists all users (admin only)
// @Summary List users
// @Description Get all users (admin only)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Server error"
// @Router /admin/users [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	err := db.DB.Find(&users).Error
	if err != nil {
		http.Error(w, "error with fetching users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser updates a user (admin or self)
// @Summary Update user
// @Description Update a user's details
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user data"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/updateuser/"):]

	var user models.User 
	err := db.DB.First(&user, "id = ?", id).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if user.ID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if user.Password != "" {
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error with hashing password", http.StatusInternalServerError)
			return
		}
		user.Password = string(hashedpassword)
	}
	err = db.DB.Save(&user).Error
	if err != nil {
		http.Error(w, "Error with updating", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/deleteusers/"):]
	var user models.User
	err := db.DB.First(&user, "id = ?", id).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	err = db.DB.Delete(&user).Error
	if err != nil {
		http.Error(w, "Error with deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

