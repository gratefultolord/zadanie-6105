package models

import (
	"time"
)

type OrganizationType string

const (
	OrganizationTypeIE  OrganizationType = "IE"
	OrganizationTypeLLC OrganizationType = "LLC"
	OrganizationTypeJSC OrganizationType = "JSC"
)

type Organization struct {
	ID          string           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string           `gorm:"type:varchar(100);not null" json:"name"`
	Description string           `gorm:"type:text" json:"description"`
	Type        OrganizationType `gorm:"type:organization_type" json:"type"`
	CreatedAt   time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
}
