package repository

import "kaabe-app/internal/domain/model"

type TokenRepository interface {
	FindByToken(token string) (*model.Token, error)
	Create(token *model.Token) error
}
