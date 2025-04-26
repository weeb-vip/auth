package user_client

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/weeb-vip/auth/config"
)

type UserClientInterface interface {
	AddUser(ctx context.Context, id string, username string, firstname string, lastname string, language string) error
}

type userClient struct {
	client *graphql.Client
}

func NewUserClient(cfg config.UserClientConfig) UserClientInterface {
	client := graphql.NewClient(cfg.URL)
	return &userClient{
		client: client,
	}
}

func (u *userClient) AddUser(ctx context.Context, id, username, firstname, lastname, language string) error {
	req := graphql.NewRequest(`
		mutation AddUser($input: AddUserInput!) {
			AddUser(input: $input) {
				id
			}
		}
	`)

	req.Var("input", map[string]interface{}{
		"id":        id,
		"firstname": firstname,
		"lastname":  lastname,
		"username":  username,
		"language":  language,
	})

	req.Header.Set("Cache-Control", "no-cache")

	// Define a strongly typed response structure
	type AddUserResponse struct {
		AddUser struct {
			ID string `json:"id"`
		} `json:"AddUser"`
	}

	var respData AddUserResponse
	err := u.client.Run(ctx, req, &respData)
	if err != nil {
		return fmt.Errorf("failed to execute GraphQL request: %w", err)
	}

	// Check if the response contains valid data
	if respData.AddUser.ID == "" {
		return fmt.Errorf("failed to add user: missing user ID in response")
	}

	return nil
}
