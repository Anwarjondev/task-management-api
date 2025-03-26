package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string `gorm:"type:varchar(255)" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	OwnerID     string `gorm:"type:uuid" json:"owner_id"`
	Owner       User   `gorm:"foreignKey:OwnerID" json:"owner"`
	Members     []User `gorm:"many2many:project_members;" json:"members"`
	Tasks       []Task `gorm:"foreignKey:ProjectID" json:"tasks"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	return nil
}
