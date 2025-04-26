package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/weeb-vip/auth/internal/db"
	"github.com/weeb-vip/auth/internal/services/session/models"
)

type SessionsRepository interface {
	CreateSession(ctx context.Context, userID string) (*models.Session, error)
	GetSession(ctx context.Context, token string) (*models.Session, error)
	DeleteSession(ctx context.Context, token string) error
}

type sessionsRepository struct {
	DBService db.DB
}

var sessionsRepositorySingleton SessionsRepository // nolint

func NewSessionsRepository() SessionsRepository {
	dbService := db.GetDBService()

	return &sessionsRepository{
		DBService: dbService,
	}
}

func (repository *sessionsRepository) CreateSession(
	ctx context.Context,
	userID string,
) (*models.Session, error) {
	database := repository.DBService.GetDB()

	session := models.Session{
		UserID:    userID,
		IPAddress: "",
		UserAgent: "",
		Token:     "",
	}

	err := database.Create(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (repository *sessionsRepository) GetSession(
	ctx context.Context,
	token string,
) (*models.Session, error) {
	database := repository.DBService.GetDB()

	var session models.Session

	err := database.Where("token = ?", token).First(&session).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &session, nil
}

func (repository *sessionsRepository) DeleteSession(
	ctx context.Context,
	token string,
) error {
	database := repository.DBService.GetDB()

	return database.Where("token = ?", token).Delete(&models.Session{}).Error
}

func GetSessionsRepository() SessionsRepository {
	if sessionsRepositorySingleton == nil {
		sessionsRepositorySingleton = NewSessionsRepository()
	}

	return sessionsRepositorySingleton
}
