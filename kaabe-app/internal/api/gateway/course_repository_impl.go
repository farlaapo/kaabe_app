package gateway

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

type CourseRepositoryImpl struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) repository.CourseRepository {
	return &CourseRepositoryImpl{db: db}
}

// Create inserts a new course using the stored procedure
func (r *CourseRepositoryImpl) Create(course *model.Course) error {
	_, err := r.db.Exec(`CALL create_course($1, $2, $3, $4, $5, $6, $7)`,
		course.ID, // Include course.ID here
		course.InfluencerID,
		course.Title,
		course.Description,
		course.Price,
		pq.Array(course.CoverImageURL),
		course.Status,
	)

	if err != nil {
		log.Printf("Error calling create_course: %v", err)
		return err
	}

	log.Printf("Course created: %+v", course)
	return nil
}

// Update modifies an existing course using the stored procedure
func (r *CourseRepositoryImpl) Update(course *model.Course) error {
	_, err := r.db.Exec(`CALL update_course($1, $2, $3, $4, $5, $6)`,
		course.ID, course.Title, course.Description, course.Price, pq.Array(course.CoverImageURL), course.Status)
	if err != nil {
		log.Printf("Error calling update_course: %v", err)
		return err
	}

	log.Printf("Course updated: %+v", course)
	return nil
}

// Delete performs a soft delete of a course using the stored procedure
func (r *CourseRepositoryImpl) Delete(courseID uuid.UUID) error {
	_, err := r.db.Exec(`CALL delete_course($1)`, courseID)
	if err != nil {
		log.Printf("Error calling delete_course for ID %v: %v", courseID, err)
		return err
	}

	log.Printf("Course soft-deleted: %v", courseID)
	return nil
}

// GetByID retrieves a single course by its ID using the get_course_by_id() function
func (r *CourseRepositoryImpl) GetByID(courseID uuid.UUID) (*model.Course, error) {
	var course model.Course

	query := `SELECT * FROM get_course_by_id($1)`
	row := r.db.QueryRow(query, courseID)

	err := row.Scan(
		&course.ID,
		&course.InfluencerID,
		&course.Title,
		&course.Description,
		&course.Price,
		pq.Array(&course.CoverImageURL),
		&course.Status,
		&course.CreatedAt,
		&course.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("course not found")
			return nil, fmt.Errorf("course not found")
		}
		log.Printf("DB error: %v", err)
		return nil, err
	}

	return &course, nil
}

// GetAll retrieves all non-deleted courses using the get_all_courses() function
func (r *CourseRepositoryImpl) GetAll() ([]*model.Course, error) {
	rows, err := r.db.Query(`SELECT * FROM get_all_courses()`)
	if err != nil {
		log.Printf("Error querying get_all_courses: %v", err)
		return nil, err
	}
	defer rows.Close()

	var courses []*model.Course

	for rows.Next() {
		var course model.Course
		err := rows.Scan(
			&course.ID,
			&course.InfluencerID,
			&course.Title,
			&course.Description,
			&course.Price,
			pq.Array(&course.CoverImageURL),
			&course.Status,
			&course.CreatedAt,
			&course.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning course row: %v", err)
			return nil, err
		}
		courses = append(courses, &course)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("Courses retrieved: %d", len(courses))
	return courses, nil
}
