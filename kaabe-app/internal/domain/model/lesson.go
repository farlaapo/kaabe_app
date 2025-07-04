package model

import (
	"time"

	"github.com/gofrs/uuid"
)

type Lesson struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CourseID  uuid.UUID `json:"course_id" gorm:"type:uuid;not null"`
	Title     string    `json:"title"`
	VideoURL  []string  `json:"video_url" gorm:"type:jsonb"`
	Order     int       `json:"order" gorm:"column:lesson_order;type:int;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
