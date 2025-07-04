package repository

import (
	"kaabe-app/internal/domain/model"

	"github.com/gofrs/uuid"
)

type LessonRepository interface {
	Create(lesson *model.Lesson) error
	Update(lesson *model.Lesson) error
	Delete(lessonID uuid.UUID) error
	GetByID(lessonID uuid.UUID) (*model.Lesson, error)
	GetAll() ([]*model.Lesson, error)
}
