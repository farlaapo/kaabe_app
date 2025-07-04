package service

import (
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"

	"log"
	"time"

	"github.com/gofrs/uuid"
)

// CourseService interface
type CourseService interface {
	CreateCourse(InfluencerID uuid.UUID, Title, Description string, Price float64, CoverImageURL []string, Status string) (*model.Course, error)
	UpdateCourse(course *model.Course) error
	DeleteCourse(courseID uuid.UUID) error
	GetCourseByID(courseID uuid.UUID) (*model.Course, error)
	GetAllCourses() ([]*model.Course, error)
}

// courseServiceImpl struct implementing CourseService
type courseServiceImpl struct {
	repo      repository.CourseRepository
	tokenRepo repository.TokenRepository
}

// CreateCourse implements CourseService.
func (c *courseServiceImpl) CreateCourse(InfluencerID uuid.UUID, Title string, Description string, Price float64 , CoverImageURL []string, Status string) (*model.Course, error) {
	// Generate a new UUID for the course ID
	neoCourse, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Create a new course instance
	amCourse := &model.Course{
		ID:            neoCourse,
		InfluencerID:  InfluencerID,
		Title:         Title,
		Description:   Description,
		Price:         Price,
		CoverImageURL: CoverImageURL,
		Status:        Status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Log the new course creation attempt
	log.Printf("Creating course: %+v", amCourse)

	// Save the new course to the repository
	err = c.repo.Create(amCourse)
	if err != nil {
		return nil, fmt.Errorf("failed to create course: %v", err)
	}

	return amCourse, nil
}

// DeleteCourse implements CourseService.
func (c *courseServiceImpl) DeleteCourse(courseID uuid.UUID) error {
	_, err := c.repo.GetByID(courseID)
	if err != nil {
		return fmt.Errorf("could not find course with ID %s: %v", courseID, err)
	}

	if err := c.repo.Delete(courseID); err != nil {
		return fmt.Errorf("failed to delete course with ID %s: %v", courseID, err)
	}

	log.Printf("Successfully deleted course with ID %s", courseID)
	return nil
}

// GetAllCourses implements CourseService.
func (c *courseServiceImpl) GetAllCourses() ([]*model.Course, error) {
	course, err := c.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all courses: %v", err)
	}
	return course, nil
}

// GetCourseByID implements CourseService.
func (c *courseServiceImpl) GetCourseByID(courseID uuid.UUID) (*model.Course, error) {
	course, err := c.repo.GetByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("could not find course with ID %s: %v", courseID, err)
	}
	return course, nil
}

// UpdateCourse implements CourseService.
func (c *courseServiceImpl) UpdateCourse(course *model.Course) error {
	_, err := c.repo.GetByID(course.ID)
	if err != nil {
		return fmt.Errorf("could not find course with ID %s", course.ID)
	}

	if err := c.repo.Update(course); err != nil {
		return fmt.Errorf("failed to update course with ID %s: %v", course.ID, err)
	}

	return nil
}

// NewCourseService creates a new instance of CourseService
func NewCourseService(coureRepo repository.CourseRepository, tokenRepo repository.TokenRepository) CourseService {
	return &courseServiceImpl{
		repo:      coureRepo,
		tokenRepo: tokenRepo,
	}
}
