package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Course struct {
    ID            uuid.UUID `json:"id,omitempty" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    InfluencerID  uuid.UUID `json:"influencer_id,omitempty" gorm:"type:uuid;not null"`
    Title         string    `json:"title" gorm:"type:varchar(255);not null"`
    Description   string    `json:"description" gorm:"type:text"`
    Price         float64   `json:"price" gorm:"type:float"`
    CoverImageURL []string  `json:"cover_image_url" gorm:"type:text[]"`  // slice of strings, matches Postgres TEXT[]
    Status        string    `json:"status" gorm:"type:varchar(50)"`
    CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    
}

