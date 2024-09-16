package models

import (
	"time"
)

type TenderStatus string

const (
	TenderStatusCreated   TenderStatus = "Created"
	TenderStatusPublished TenderStatus = "Published"
	TenderStatusClosed    TenderStatus = "Closed"
)

type TenderServiceType string

const (
	ServiceTypeConstruction TenderServiceType = "Construction"
	ServiceTypeDelivery     TenderServiceType = "Delivery"
	ServiceTypeManufacture  TenderServiceType = "Manufacture"
)

type Tender struct {
	ID              string            `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Version         int               `gorm:"primaryKey" json:"version"`
	Name            string            `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Description     string            `gorm:"type:text;not null" json:"description" validate:"required"`
	ServiceType     TenderServiceType `gorm:"type:varchar(50);not null" json:"serviceType" validate:"required"`
	OrganizationID  string            `gorm:"type:uuid;not null" json:"organizationId" validate:"required,uuid4"`
	CreatorUsername string            `gorm:"type:varchar(50);not null" json:"creatorUsername" validate:"required"`
	Status          TenderStatus      `gorm:"type:varchar(50);default:'Created'" json:"status"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"createdAt"`
}
