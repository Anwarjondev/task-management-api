package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
)



// CreateProject creates a new project
// @Summary Create a project
// @Description Create a new project owned by the authenticated user
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body models.Project true "Project data"
// @Success 201 {object} models.Project
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /createproject [post]
func CreateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	project.OwnerID = userID
	err = db.DB.Create(&project).Error
	if err != nil {
		http.Error(w, "Error: Creating project", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}


// GetProjects lists projects with pagination
// @Summary List projects
// @Description Get projects accessible to the user with pagination
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Success 200 {array} models.Project
// @Failure 401 {string} string "Unauthorized"
// @Router /getproject [get]
func GetProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	userRole := r.Context().Value("user_role").(string)
	var projects []models.Project
	if userRole == "admin" {
		db.DB.Find(&projects)
	} else {
		db.DB.Joins("Join project_members on project_members.project_id = project.id").Where("project_members.user_id = ? or project.owner_id = ?", userID, userID).Find(&projects)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// UpdateProject updates a project
// @Summary Update a project
// @Description Update a project if the user is the owner or admin
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param project body models.Project true "Updated project data"
// @Success 200 {object} models.Project
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Server error"
// @Router /projects/{id} [put]
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/updateproject/"):]

	var project models.Project

	err := db.DB.First(&project, "id = ?", id).Error
	if err != nil {
		http.Error(w, "Project is not found", http.StatusNotFound)
		return
	}
	if project.OwnerID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = db.DB.Save(&project).Error
	if err != nil {
		http.Error(w, "Error updating project", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// DeleteProject deletes a project
// @Summary Delete a project
// @Description Delete a project if the user is the owner or admin
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 204 {string} string "No content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /projects/{id} [delete]
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/deleteproject/"):]

	var project models.Project
	err := db.DB.First(&project, "id = ?", id).Error
	if err != nil {
		http.Error(w, "Project not Found", http.StatusNotFound)
		return
	}
	if project.OwnerID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = db.DB.Delete(&project).Error
	if err != nil {
		http.Error(w, "Error deleting project", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// AddProjectMember adds a user to a project
// @Summary Add project member
// @Description Add a user to a project if the user is the owner or admin
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param user body map[string]string true "User ID to add"
// @Success 200 {object} models.Project
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /projects/{id}/members [post]
func AddProjectMember(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/projects/"):]

	var project models.Project
	err := db.DB.Preload("Members").First(&project, "id = ?", id).Error
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}
	if project.OwnerID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var input struct {
		UserID string `json:"user_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var user models.User
	err = db.DB.First(&user, "id = ?", input.UserID).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	project.Members = append(project.Members, user)
	err = db.DB.Save(&project).Error
	if err != nil {
		http.Error(w, "Error adding members", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}