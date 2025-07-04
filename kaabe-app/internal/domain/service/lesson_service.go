package service

import (
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"

	"log"
	"time"

	"github.com/gofrs/uuid"
)

type LessonService interface {
	CreateLesson(CourseID uuid.UUID, title string, VideoURL []string, Order int) (*model.Lesson, error)
	UpdateLesson(lesson *model.Lesson) error
	DeleteLesson(LessonID uuid.UUID) error
	GetLessonByID(lessonID uuid.UUID) (*model.Lesson, error) 
	GetAllLessons() ([]*model.Lesson, error)
}

// lessonServiceImpl struct implementing lessonService
type lessonServiceImpl struct {
	repo      repository.LessonRepository
	tokenRepo repository.TokenRepository
}

func NewLessonService(lessonRepo repository.LessonRepository, tokenRepo repository.TokenRepository) LessonService {
	return &lessonServiceImpl{
		repo:      lessonRepo,
		tokenRepo: tokenRepo,
	}
}

// CreateLesson implements LessonService.
func (l *lessonServiceImpl) CreateLesson(CourseID uuid.UUID, title string, VideoURL []string, Order int) (*model.Lesson, error) {
	// Generate a new UUID for the lesson ID

	neoLesson, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Create a new lesson instance
	amLesson := &model.Lesson{
		ID:        neoLesson,
		CourseID:  CourseID,
		Title:     title,
		VideoURL:  VideoURL,
		Order:     Order,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Log the new lesson creation attempt
	log.Printf("Creating lesson: %+v", amLesson)

	// Save the new lesson to the repository
	err = l.repo.Create(amLesson)
	if err != nil {
		return nil, fmt.Errorf("failed to create lesson: %v", err)
	}

	return amLesson, nil
}

// DeleteLesson implements LessonService.
func (l *lessonServiceImpl) DeleteLesson(LessonID uuid.UUID) error {
	_, err := l.repo.GetByID(LessonID)
	if err != nil {
		return fmt.Errorf("could not find lesson with ID %s: %v", LessonID, err)
	}

	if err := l.repo.Delete(LessonID); err != nil {
		return fmt.Errorf("failed to delete lesson with ID %s: %v", LessonID, err)
	}

	log.Printf("Successfully deleted lesson with ID %s", LessonID)
	return nil
}

// GetAllLessons implements LessonService.
func (l *lessonServiceImpl) GetAllLessons() ([]*model.Lesson, error) {
	lesson, err := l.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all lesson: %v", err)
	}
	return lesson, nil
}

// GetLessonByID retrieves a lesson by its ID.
func (l *lessonServiceImpl) GetLessonByID(lessonID uuid.UUID) (*model.Lesson, error) {
	lesson, err := l.repo.GetByID(lessonID)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

// UpdateLesson implements LessonService.
func (l *lessonServiceImpl) UpdateLesson(lesson *model.Lesson) error {
	_, err := l.repo.GetByID(lesson.ID)
	if err != nil {
		return fmt.Errorf("could not find lesson with ID %s", lesson.ID)
	}

	if err := l.repo.Update(lesson); err != nil {
		return fmt.Errorf("failed to update lesson with ID %s: %v", lesson.ID, err)
	}

	return nil
}
