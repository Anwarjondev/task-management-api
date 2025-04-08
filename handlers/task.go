package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Anwarjondev/task-management-api/db"
	"github.com/Anwarjondev/task-management-api/models"
)

// CreateTask creates a new task
// @Summary Create a task
// @Description Create a task in a project
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.Task true "Task data"
// @Success 201 {object} models.Task
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /createtask [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	task.CreatorID = userID
	task.Status = "pending"
	err = db.DB.Create(&task).Error
	if err != nil {
		http.Error(w, "Error with creating task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTasks lists tasks with pagination
// @Summary List tasks
// @Description Get tasks accessible to the user with pagination
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {array} models.Task
// @Failure 401 {string} string "Unauthorized"
// @Router /gettask [get]
func GetTask(w http.ResponseWriter, r *http.Request) {
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
	offset := (page-1) * perPage
	status := r.URL.Query().Get("status")

	var tasks []models.Task
	query := db.DB.Limit(perPage).Offset(offset)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if role == "admin" {
		query.Find(&tasks)
	} else {
		query.Where("creator_id = ? or assignee_id = ?", userID, userID).Find(&tasks)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}


// UpdateTask updates a task
// @Summary Update a task
// @Description Update a task if the user is the creator, assignee, or admin
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param task body models.Task true "Updated task data"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /tasks/{id} [put]
func Updatetask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	role := r.Context().Value("role").(string)
	id := r.URL.Path[len("/updatetask/"):]
	
	var task models.Task
	err := db.DB.First(&task, "id = ?", id).Error
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if task.AssigneeID != userID && task.CreatorID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err = db.DB.Save(&task).Error
	if err != nil {
		http.Error(w, "Error with supdating task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask deletes a task
// @Summary Delete a task
// @Description Delete a task if the user is the creator or admin
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 204 {string} string "No content"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not found"
// @Router /tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	role := r.Context().Value("role")
	id := r.URL.Path[len("/deletetask/"):]

	var task models.Task
	err := db.DB.First(&task, "id = ?", id).Error
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	if task.CreatorID != userID && role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err = db.DB.Delete(&task).Error 
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


