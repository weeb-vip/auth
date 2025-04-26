package resolvers

import (
	"context"
	"errors"
	"log"

	"github.com/99designs/gqlgen/graphql"

	"github.com/weeb-vip/auth/internal/entities"
	"github.com/weeb-vip/auth/internal/services/credential"
	"github.com/weeb-vip/auth/internal/xerrors"
)

func handleError(ctx context.Context, result interface{}, err error) (interface{}, error) { //nolint
	log.Println(result)
	log.Println(err)

	graphql.AddError(ctx, xerrors.ServiceError(err.Error(), getCode(err)))

	return result, nil
}

func getCode(err error) string {
	var credErr *credential.Error
	if ok := errors.As(err, &credErr); ok {
		return credErr.Code.String()
	}

	var servErr *entities.ServiceError
	if ok := errors.As(err, &servErr); ok {
		return servErr.Code
	}

	return "UNKNOWN_ERROR"
}
