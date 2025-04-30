package resolvers

import (
	"context"
	"log"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/internal/services/validation_token"

	"github.com/weeb-vip/auth/internal/services/credential"
)

func Register( // nolint
	ctx context.Context,
	authenticationService credential.Credential,
	validatonToken validation_token.ValidationToken,
	conf config.Config,
	firstName string,
	lastName string,
	username string,
	password string,
	language string,
) (*model.RegisterResult, error) {
	credentials, err := authenticationService.Register(ctx, username, password)
	if err != nil {
		res, err := handleError(ctx, "null", err)
		if res != nil {
			return nil, err
		}

		return nil, err
	}

	token, err := validatonToken.GenerateToken(username)
	if err != nil {
		return nil, err
	}

	//err = notificationService.Publish(conf.SNSConfig.RegistrationCompleteTopic, types.RegistrationComplete{
	//	EventName:         "registration_complete",
	//	FirstName:         firstName,
	//	LastName:          lastName,
	//	Identifier:        credentials.Username,
	//	IdentifierType:    types.Email,
	//	PreferredLanguage: types.PreferredLanguage(language),
	//	RegisteredAt:      credentials.CreatedAt.Format("2006-01-02 15:04:05"),
	//	SchemaVersion:     "1.0",
	//	UserID:            credentials.UserID,
	//	VerificationToken: token,
	//})
	//
	//if err != nil {
	//	return nil, err
	//}

	log.Println(token)

	return &model.RegisterResult{
		ID: credentials.UserID,
	}, nil
}
