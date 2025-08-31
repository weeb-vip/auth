package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/weeb-vip/auth/internal/db"
	"github.com/weeb-vip/auth/internal/services/credential/models"
)

type CredentialsRepository interface {
	AddCredentials(
		username string,
		userID string,
		value string,
		credType models.CredentialTypes,
	) (*models.Credential, error)
	GetCredentials(username string) (*models.Credential, error)
	DeleteCredentials(username string) error
	UpdatePassword(username string, hashedPassword string) error
	ActivateCredentials(id string) error
}

type credentialsRepository struct {
	DBService db.DB
}

var credentialsRepositorySingleton CredentialsRepository // nolint

func NewCredentialsRepository() CredentialsRepository {
	dbService := db.GetDBService()

	return &credentialsRepository{
		DBService: dbService,
	}
}

func (repository *credentialsRepository) AddCredentials(
	username string,
	userID string,
	value string,
	credType models.CredentialTypes,
) (*models.Credential, error) {
	database := repository.DBService.GetDB()

	credentials := models.Credential{
		Username: username,
		UserID:   userID,
		Value:    value,
		Type:     credType,
	}
	err := database.FirstOrCreate(&credentials, models.Credential{Username: username}).Error

	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

func (repository *credentialsRepository) GetCredentials(username string) (*models.Credential, error) {
	database := repository.DBService.GetDB()

	var credentials models.Credential

	err := database.Where("username = ?", username).First(&credentials).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &credentials, nil
}

func (repository *credentialsRepository) DeleteCredentials(username string) error {
	database := repository.DBService.GetDB()

	return database.Where("username = ?", username).Delete(&models.Credential{}).Error
}

func (repository *credentialsRepository) UpdatePassword(username string, hashedPassword string) error {
	database := repository.DBService.GetDB()

	return database.Model(&models.Credential{}).
		Where("username = ?", username).
		Update("value", hashedPassword).Error
}

func GetCredentialsRepository() CredentialsRepository {
	if credentialsRepositorySingleton == nil {
		credentialsRepositorySingleton = NewCredentialsRepository()
	}

	return credentialsRepositorySingleton
}

func (repository *credentialsRepository) ActivateCredentials(id string) error {
	database := repository.DBService.GetDB()

	return database.Model(&models.Credential{}).
		Where("id = ?", id).
		Update("active", true).Error
}
