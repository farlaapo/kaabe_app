package gateway

import (
	"database/sql"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

type LessonRepositoryImpl struct {
	db *sql.DB
}

// Create inserts a new lesson using the stored procedure
func (l *LessonRepositoryImpl) Create(lesson *model.Lesson) error {
	_, err := l.db.Exec(`CALL create_lesson($1, $2, $3, $4, $5)`,
		lesson.ID, lesson.CourseID, lesson.Title, pq.Array(lesson.VideoURL), lesson.Order)
	if err != nil {
		log.Printf("Error calling create_lesson: %v", err)
		return err
	}

	log.Printf("Lesson created: %+v", lesson)
	return nil
}



// Delete implements repository.LessonRepository.
// Delete performs a soft delete of a lesson using the stored procedure
func (l *LessonRepositoryImpl) Delete(lessonID uuid.UUID) error {
	_, err := l.db.Exec(`CALL  delete_lesson($1)`, lessonID)
	if err != nil {
		log.Printf("Error calling delete_lesson for ID %v: %v", lessonID, err)
		return err
	}

	log.Printf("Lesson soft-deleted: %v", lessonID)
	return nil
}


// GetAllByCourseID retrieves all lessons for a specific course using the get_lessons_by_course_id() function
func (l *LessonRepositoryImpl) GetAll() ([]*model.Lesson, error) {
	rows, err := l.db.Query(`SELECT * FROM get_all_lessons()`)
	if err != nil {
		log.Printf("Error querying get_all_lessons: %v", err)
		return nil, err
	}
	defer rows.Close()

	var lessons []*model.Lesson

	for rows.Next() {
		var lesson model.Lesson
		err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Title,
			pq.Array(&lesson.VideoURL),
			&lesson.Order,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning lesson row: %v", err)
			return nil, err
		}
		lessons = append(lessons, &lesson)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("Lessons retrieved: %d", len(lessons))
	return lessons, nil
}


// GetByID implements repository.LessonRepository.
func (l *LessonRepositoryImpl) GetByID(lessonID uuid.UUID) (*model.Lesson, error) {
	var lesson model.Lesson

	row := l.db.QueryRow(`SELECT * FROM get_lesson_by_id($1)`, lessonID)

	err := row.Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Title,
		pq.Array(&lesson.VideoURL),
		&lesson.Order,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Lesson not found with ID: %v", lessonID)
			return nil, fmt.Errorf("lesson not found")
		}
		log.Printf("Error scanning lesson by ID: %v", err)
		return nil, err
	}

	log.Printf("Lesson retrieved by ID: %+v", lesson)
	return &lesson, nil
}

// Update implements repository.LessonRepository.
// Update modifies an existing lesson using the stored procedure
func (l *LessonRepositoryImpl) Update(lesson *model.Lesson) error {
	_, err := l.db.Exec(`CALL update_lesson($1, $2, $3, $4)`,
		lesson.ID, lesson.Title, pq.Array(lesson.VideoURL), lesson.Order)
	if err != nil {
		log.Printf("Error calling update_lesson: %v", err)
		return err
	}

	log.Printf("Lesson updated: %+v", lesson)
	return nil
}

func NewLessonRepository(db *sql.DB) repository.LessonRepository {
	return &LessonRepositoryImpl{db: db}
}
