package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/weeb-vip/auth/internal/db"
	"github.com/weeb-vip/auth/internal/services/refresh_token/models"
)

type RefreshTokenRepository interface {
	AddRefreshToken(userID string, token string, expiry int64) (*models.RefreshToken, error)
	GetRefreshToken(token string) (*models.RefreshToken, error)
	GetRefreshTokenByUserID(userID string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
}

type refreshTokenRepository struct {
	DBService db.DB
}

var refreshTokenRepositorySingleton RefreshTokenRepository // nolint

func NewRefreshTokenRepository() RefreshTokenRepository {
	dbService := db.GetDBService()

	return &refreshTokenRepository{
		DBService: dbService,
	}
}

func (repository *refreshTokenRepository) AddRefreshToken(
	userID string,
	token string,
	expiry int64,
) (*models.RefreshToken, error) {
	database := repository.DBService.GetDB()

	refreshToken := models.RefreshToken{
		UserID: userID,
		Token:  token,
		Expiry: expiry,
	}

	err := database.Create(&refreshToken).Error
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (repository *refreshTokenRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	database := repository.DBService.GetDB()

	var refreshToken models.RefreshToken

	err := database.Where("token = ?", token).First(&refreshToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &refreshToken, nil
}

func (repository *refreshTokenRepository) DeleteRefreshToken(token string) error {
	database := repository.DBService.GetDB()

	err := database.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}

func GetRefreshTokenRepository() RefreshTokenRepository {
	if refreshTokenRepositorySingleton == nil {
		refreshTokenRepositorySingleton = NewRefreshTokenRepository()
	}

	return refreshTokenRepositorySingleton
}

func (repository *refreshTokenRepository) GetRefreshTokenByUserID(userID string) (*models.RefreshToken, error) {
	database := repository.DBService.GetDB()

	var refreshToken models.RefreshToken

	err := database.Where("user_id = ?", userID).First(&refreshToken).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &refreshToken, nil
}
