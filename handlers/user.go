package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
	"github.com/Anwarjondev/task-management-api/utils"
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
		utils.SendError(w, http.StatusInternalServerError, "error with fetching users: "+err.Error())
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
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /users/{id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/updateuser/"):]

	var user models.User 
	err := db.DB.First(&user, "id = ?", id).Error
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "User not found: "+err.Error())
		return
	}
	if user.ID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "Forbidden for updating user")
		return
	}
	var updateUser models.User
	err = json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	err = validate.Struct(&updateUser)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	user.Username = updateUser.Username
	if updateUser.Password != "" {
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Error with hashing password: "+err.Error())
			return
		}
		user.Password = string(hashedpassword)
	}
	user.Role = updateUser.Role
	err = db.DB.Save(&user).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error with updating: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user (admin only)
// @Summary Delete user
// @Description Delete a user (admin only)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 {string} string "No content"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Not found"
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /admin/users/{id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/deleteusers/"):]
	var user models.User
	err := db.DB.First(&user, "id = ?", id).Error
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "User not found: "+err.Error())
		return
	}
	err = db.DB.Delete(&user).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error with deleting user: "+err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

