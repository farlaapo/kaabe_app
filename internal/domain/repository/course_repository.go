package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

// CourseRepository interface with required methods
type CourseRepository interface {
	Create(course *model.Course) error
	Update(course *model.Course) error
	Delete(courseID uuid.UUID) error
	GetByID(courseID uuid.UUID) (*model.Course, error)
	GetAll() ([]*model.Course, error)
}
