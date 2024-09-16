package models

import (
	"time"
)

type Employee struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	FirstName string    `gorm:"type:varchar(50)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(50)" json:"last_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Employee) TableName() string {
	return "employee"
}
