package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title       string    `gorm:"type:varchar(255)" json:"title" validate:"required, min=3, max=100"`
	Description string    `gorm:"type:text" json:"description" validate:"max=500"`
	Status      string    `gorm:"type:varchar(50)" json:"status" validate:"required,oneof=pending in_progress completed"`
	ProjectID   string    `gorm:"type:uuid" json:"project_id" validate:"required"`
	Project     Project   `gorm:"foreignKey:ProjectID" json:"project"`
	AssigneeID  string    `gorm:"type:uuid" json:"assignee_id"`
	Assignee    User      `gorm:"foreignKey:AssigneeID" json:"assignee"`
	CreatorID   string    `gorm:"type:uuid" json:"creator_id"`
	Creator     User      `gorm:"foreignKey:CreatorID" json:"creator"`
	Subtasks    []Subtask `gorm:"foreignKey:TaskID" json:"subtasks"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New().String()
	return nil
}
