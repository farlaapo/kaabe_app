package gateway

import (
	"database/sql"
	"fmt"
	"kaabe-app/internal/domain/model"
	"kaabe-app/internal/domain/repository"
	"log"

	"github.com/gofrs/uuid"
)

type RatingRepositoryImpl struct {
	db *sql.DB
}

// Create implements repository.RatingRepository.
func (r *RatingRepositoryImpl) Create(rating *model.Rating) error {
	_, err := r.db.Exec(`CALL create_rating($1, $2, $3, $4, $5)`,
		rating.ID, rating.UserID, rating.CourseID, rating.Score, rating.Comment)
	if err != nil {
		log.Printf("Error calling create_rating: %v", err)
		return err
	}

	log.Printf("rating created: %+v", rating)
	return nil
}

// Delete implements repository.RatingRepository.
func (r *RatingRepositoryImpl) Delete(ratingID uuid.UUID) error {
	_, err := r.db.Exec(`CALL  delete_rating($1)`, ratingID)
	if err != nil {
		log.Printf("Error calling delete_rating for ID %v: %v", ratingID, err)
		return err
	}

	log.Printf("Lesson soft-deleted: %v", ratingID)
	return nil
}

// GetByID implements repository.RatingRepository.
func (r *RatingRepositoryImpl) GetByID(ratingID uuid.UUID) (*model.Rating, error) {
	var rating model.Rating

	row := r.db.QueryRow(`SELECT * FROM get_rating_by_id($1)`, ratingID)

	err := row.Scan(
		&rating.ID,
		&rating.UserID,
		&rating.CourseID,
		&rating.Score,
		&rating.Comment,
		&rating.CreatedAt,
		&rating.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("rating not found with ID: %v", ratingID)
			return nil, fmt.Errorf("rating not found")
		}
		log.Printf("Error scanning rating by ID: %v", err)
		return nil, err
	}

	log.Printf("rating retrieved by ID: %+v", rating)
	return &rating, nil
}

// Getall implements repository.RatingRepository.
func (r *RatingRepositoryImpl) Getall() ([]*model.Rating, error) {

	rows, err := r.db.Query(`SELECT * FROM get_all_ratings()`)
	if err != nil {
		log.Printf("Error querying get_all_ratings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var Ratings []*model.Rating

	for rows.Next() {
		var Rating model.Rating
		err := rows.Scan(
			&Rating.ID,
			&Rating.UserID,
			&Rating.CourseID,
			&Rating.Score,
			&Rating.Comment,
			&Rating.CreatedAt,
			&Rating.UpdatedAt,
		)

		if err != nil {
			log.Printf("Error scanning ratig row: %v", err)
			return nil, err
		}
		Ratings = append(Ratings, &Rating)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	log.Printf("Lessons retrieved: %d", len(Ratings))
	return Ratings, nil
}

// Update implements repository.RatingRepository.
func (r *RatingRepositoryImpl) Update(rating *model.Rating) error {
	_, err := r.db.Exec(`CALL update_rating($1, $2, $3)`,
		rating.ID, rating.Score, rating.Comment)
	if err != nil {
		log.Printf("Error calling update_rating: %v", err)
		return err
	}

	log.Printf("Lesson updated: %+v", rating)
	return nil
}

func NewRatingRepository(db *sql.DB) repository.RatingRepository {
	return &RatingRepositoryImpl{db: db}
}
