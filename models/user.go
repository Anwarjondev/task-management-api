package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();unique" json:"id"`
	Username string    `gorm:"type:varchar(255);unique" json:"username"`
	Password string    `gorm:"type:varchar(255)" json:"password"`
	Role     string    `gorm:"type:varchar(50)" json:"role"`
	Projects []Project `gorm:"many2many:project_members;" json:"projects"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}
