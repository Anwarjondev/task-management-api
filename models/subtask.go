package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subtask struct {
	ID         string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Title      string `gorm:"type:varchar(255)" json:"title" validate:"required,min=3, max=100"`
	Status     string `gorm:"type:varchar(50)" json:"status" validate:"required, oneof=pending, in_progress completed"`
	TaskID     string `gorm:"type uuid" json:"task_id" validate:"required"`
	Task       Task   `gorm:"foreignKey:TaskID" json:"task"`
	AssigneeID string `gorm:"type:uuid" json:"assignee_id"`
	Assignee   User   `gorm:"foreignKey:AssigneeID" json:"assignee"`
	CreatorID  string `gorm:"type:uuid" json:"creator_id"`
	Creator    User   `gorm:"foreignKey:CreatorID" json:"creator"`
}


func (s *Subtask) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New().String()
	return nil
}
