package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
	"github.com/Anwarjondev/task-management-api/utils"
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
		utils.SendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	err = validate.Struct(&project)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	project.OwnerID = userID
	err = db.DB.Create(&project).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error with creating project: "+err.Error())
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
	role := r.Context().Value("role").(string)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 10
	}
	offset := (page -1) *perPage
	var projects []models.Project

	query := db.DB.Limit(perPage).Offset(offset)
	if role == "admin" {
		err := query.Find(&projects).Error
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Error fetch projects: "+err.Error())
			return
		}
		
	} else {
		err := query.Joins("Join project_members on project_members.project_id = project.id").Where("project_members.user_id = ? or project.owner_id = ?", userID, userID).Find(&projects).Error
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch all projects: "+err.Error())
			return
		}
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
		utils.SendError(w, http.StatusNotFound, "Project Not found: "+err.Error())
		return
	}
	if project.OwnerID != userID && role != "admin" {
		utils.SendError(w,http.StatusNotFound, "You are not authorized: "+err.Error())
		return
	}
	var updateProject models.Project
	err = json.NewDecoder(r.Body).Decode(&updateProject)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Bad request: "+err.Error())
		return
	}
	err = validate.Struct(&updateProject)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	project.Name = updateProject.Name
	project.Description = updateProject.Description
	err = db.DB.Save(&project).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update project: "+err.Error())
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
		utils.SendError(w, http.StatusNotFound, "Project not found: "+err.Error())
		return
	}
	if project.OwnerID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "You are not authorized to delete: "+err.Error())
		return
	}
	err = db.DB.Delete(&project).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error deleting project: "+err.Error())
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
	id := r.URL.Path[len("/projects/"):len(r.URL.Path)-len("/members")]

	var project models.Project
	err := db.DB.Preload("Members").First(&project, "id = ?", id).Error
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Project not found: "+err.Error())
		return
	}
	if project.OwnerID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "Forbidden: "+err.Error())
		return
	}
	var input struct {
		UserID string `json:"user_id" validate:"required"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
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
