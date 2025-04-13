package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
	"github.com/Anwarjondev/task-management-api/utils"
)

// CreateSubtask creates a new subtask
// @Summary Create a subtask
// @Description Create a subtask under a task
// @Tags Subtasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subtask body models.Subtask true "Subtask data"
// @Success 201 {object} models.Subtask
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /subtasks [post]
func CreateSubTask(w http.ResponseWriter, r *http.Request) {	
	userID := r.Context().Value("user_id").(string)
	var subtask models.Subtask
	err := json.NewDecoder(r.Body).Decode(&subtask)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	err = validate.Struct(&subtask)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	subtask.CreatorID = userID
	subtask.Status = "pending"
	err = db.DB.Create(&subtask).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error with creating subtask: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subtask)
}



// GetSubtasks lists subtasks with pagination
// @Summary List subtasks
// @Description Get subtasks accessible to the user with pagination
// @Tags Subtasks
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param task_id query string false "Filter by task ID"
// @Success 200 {array} models.Subtask
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /subtasks [get]
func GetSubtask(w http.ResponseWriter, r *http.Request) {
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
	offset := (page - 1) * perPage
	taskID := r.URL.Query().Get("task_id")

	var subtasks []models.Subtask
	query := db.DB.Limit(perPage).Offset(offset)
	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if role == "admin" {
		err := query.Find(&subtasks).Error
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Error fetching subtasks: "+err.Error())
			return
		}
	} else {
		if err := query.Where("creator_id = ? or assignee_id = ?", userID, userID).Find(&subtasks).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Error fetching subtasks: "+err.Error())
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subtasks)
}

// UpdateSubtask updates a subtask
// @Summary Update a subtask
// @Description Update a subtask if the user is the creator, assignee, or admin
// @Tags Subtasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subtask پردی ID"
// @Param subtask body models.Subtask true "Updated subtask data"
// @Success 200 {object} models.Subtask
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Not found"
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /subtasks/{id} [put]
func UpdateSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/subtasks/"):]

	var subtask models.Subtask
	err := db.DB.First(&subtask, "id = ?", id).Error
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Subtask not found: "+err.Error())
		return
	}
	if subtask.CreatorID != userID && subtask.AssigneeID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "Not authorized for update subtask")
		return
	}
	var updateSubtask models.Subtask
	err = json.NewDecoder(r.Body).Decode(&updateSubtask)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	err = validate.Struct(&updateSubtask)
	if err != nil {
		utils.SendError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}
	subtask.Title = updateSubtask.Title
	subtask.Status = updateSubtask.Status
	subtask.AssigneeID = updateSubtask.AssigneeID
	err = db.DB.Save(&subtask).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Failed to update subtask")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subtask)
}

// DeleteSubtask deletes a subtask
// @Summary Delete a subtask
// @Description Delete a subtask if the user is the creator or admin
// @Tags Subtasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Subtask ID"
// @Success 204 {string} string "No content"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 403 {object} utils.ErrorResponse "Forbidden"
// @Failure 404 {object} utils.ErrorResponse "Not found"
// @Failure 500 {object} utils.ErrorResponse "Server error"
// @Router /subtasks/{id} [delete]
func DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/subtasks/"):]

	var subtask models.Subtask
	err := db.DB.First(&subtask, "id = ?", id).Error
	if err != nil {
		utils.SendError(w, http.StatusNotFound, "Subtask not found for deleting: "+err.Error())
		return
	}
	if subtask.CreatorID != userID && role != "admin" {
		utils.SendError(w, http.StatusForbidden, "You are not admin for deleting subtask: "+err.Error())
		return
	}
	err = db.DB.Delete(&subtask).Error
	if err != nil {
		utils.SendError(w, http.StatusInternalServerError, "Error with deleting subtask: "+err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}