package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/weeb-vip/auth/internal/services/passwordreset/models"

	"github.com/weeb-vip/auth/internal/db"
)

type PasswordResetRepository interface {
	AddOTT(credentialID string, ott string) (*models.PasswordReset, error)
	DeleteOTT(credentialID string) error
	GetOTT(credentialID string) (*models.PasswordReset, error)
	GetOTTByToken(ott string) (*models.PasswordReset, error)
	DeleteOTTByToken(ott string) error
}

type passwordResetRepository struct {
	DBService db.DB
}

var passwordResetRepositorySingleton PasswordResetRepository // nolint

func NewPasswordResetRepository() PasswordResetRepository {
	dbService := db.GetDBService()

	return &passwordResetRepository{
		DBService: dbService,
	}
}

func (repository *passwordResetRepository) AddOTT(credentialID string, ott string) (*models.PasswordReset, error) {
	db := repository.DBService.GetDB()

	passwordReset := models.PasswordReset{
		CredentialID: credentialID,
		OTT:          ott,
	}

	err := db.FirstOrCreate(&passwordReset, models.PasswordReset{CredentialID: credentialID}).Error

	if err != nil {
		return nil, err
	}

	return &passwordReset, nil
}

func (repository *passwordResetRepository) DeleteOTT(credentialID string) error {
	db := repository.DBService.GetDB()

	err := db.Where("credential_id = ?", credentialID).Delete(&models.PasswordReset{}).Error

	if err != nil {
		return err
	}

	return nil
}

func (repository *passwordResetRepository) GetOTT(credentialID string) (*models.PasswordReset, error) {
	db := repository.DBService.GetDB()

	var passwordReset models.PasswordReset

	err := db.Where("credential_id = ?", credentialID).First(&passwordReset).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &passwordReset, nil
}

func (repository *passwordResetRepository) GetOTTByToken(ott string) (*models.PasswordReset, error) {
	db := repository.DBService.GetDB()

	var passwordReset models.PasswordReset

	err := db.Where("ott = ?", ott).First(&passwordReset).Error

	if err != nil {
		return nil, err
	}

	return &passwordReset, nil
}

func (repository *passwordResetRepository) DeleteOTTByToken(ott string) error {
	db := repository.DBService.GetDB()

	return db.Where("ott = ?", ott).Delete(&models.PasswordReset{}).Error
}

func GetPasswordResetRepository() PasswordResetRepository {
	if passwordResetRepositorySingleton == nil {
		passwordResetRepositorySingleton = NewPasswordResetRepository()
	}

	return passwordResetRepositorySingleton
}
