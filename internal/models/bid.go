package models

import (
	"time"
)

type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
)

type AuthorType string

const (
	AuthorTypeOrganization AuthorType = "Organization"
	AuthorTypeUser         AuthorType = "User"
)

type Bid struct {
	ID          string     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string     `gorm:"type:varchar(100);not null" json:"name" validate:"required"`
	Description string     `gorm:"type:text;not null" json:"description" validate:"required"`
	Status      BidStatus  `json:"status"`
	TenderID    string     `gorm:"type:uuid;not null" json:"tenderId" validate:"required,uuid4"`
	AuthorType  AuthorType `gorm:"type:varchar(50);not null" json:"authorType" validate:"required"`
	AuthorID    string     `gorm:"type:uuid;not null" json:"authorId" validate:"required,uuid4"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	Version     int        `json:"version"`
	Feedback    string     `json:"feedback"`
}

type BidReview struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	BidID     string    `gorm:"type:uuid;not null" json:"bid_id"`
	Review    string    `gorm:"type:text;not null" json:"review"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
