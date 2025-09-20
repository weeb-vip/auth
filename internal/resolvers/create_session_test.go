package resolvers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/weeb-vip/auth/config"
	"github.com/weeb-vip/auth/graph/model"
	"github.com/weeb-vip/auth/http/handlers/responsecontext"
	"github.com/weeb-vip/auth/internal/db"
	"github.com/weeb-vip/auth/internal/jwt"
	CredentialModels "github.com/weeb-vip/auth/internal/services/credential/models"
	SessionModels "github.com/weeb-vip/auth/internal/services/session/models"
	RefreshTokenModels "github.com/weeb-vip/auth/internal/services/refresh_token/models"

	"go.uber.org/mock/gomock"
)

// Mock interfaces for testing
type MockCredentialService struct {
	ctrl *gomock.Controller
}

func NewMockCredentialService(ctrl *gomock.Controller) *MockCredentialService {
	return &MockCredentialService{ctrl: ctrl}
}

func (m *MockCredentialService) Register(ctx context.Context, username string, password string) (*CredentialModels.Credential, error) {
	return &CredentialModels.Credential{UserID: "user_123"}, nil
}

func (m *MockCredentialService) SignIn(ctx context.Context, username, password string) (*CredentialModels.Credential, error) {
	return &CredentialModels.Credential{UserID: "user_123"}, nil
}

func (m *MockCredentialService) GetCredentials(ctx context.Context, username string) (*CredentialModels.Credential, error) {
	return &CredentialModels.Credential{UserID: "user_123"}, nil
}

func (m *MockCredentialService) UpdatePassword(ctx context.Context, username string, newPassword string) error {
	return nil
}

func (m *MockCredentialService) ActivateCredentials(ctx context.Context, identifier string) error {
	return nil
}

func (m *MockCredentialService) GetCredentialsByIdentifier(ctx context.Context, identifier string) (*CredentialModels.Credential, error) {
	return &CredentialModels.Credential{UserID: "user_123"}, nil
}

type MockSessionService struct {
	ctrl *gomock.Controller
}

func NewMockSessionService(ctrl *gomock.Controller) *MockSessionService {
	return &MockSessionService{ctrl: ctrl}
}

func (m *MockSessionService) CreateSession(ctx context.Context, userID string) (*SessionModels.Session, error) {
	return &SessionModels.Session{
		BaseModel: db.BaseModel{ID: "session_123"},
		UserID:    userID,
	}, nil
}

type MockRefreshTokenService struct {
	ctrl *gomock.Controller
}

func NewMockRefreshTokenService(ctrl *gomock.Controller) *MockRefreshTokenService {
	return &MockRefreshTokenService{ctrl: ctrl}
}

func (m *MockRefreshTokenService) GetToken(token string) (*RefreshTokenModels.RefreshToken, error) {
	return &RefreshTokenModels.RefreshToken{
		Token:  token,
		UserID: "user_123",
	}, nil
}

func (m *MockRefreshTokenService) GetTokenByUserID(userID string) (*RefreshTokenModels.RefreshToken, error) {
	return nil, nil
}

func (m *MockRefreshTokenService) CreateToken(userID string) (*RefreshTokenModels.RefreshToken, error) {
	return &RefreshTokenModels.RefreshToken{
		Token:  "refresh_token_123",
		UserID: userID,
	}, nil
}

func (m *MockRefreshTokenService) DeleteToken(userID string) error {
	return nil
}

func (m *MockRefreshTokenService) ValidateToken(token string) (bool, error) {
	return true, nil
}

type MockJWTTokenizer struct {
	ctrl *gomock.Controller
}

func NewMockJWTTokenizer(ctrl *gomock.Controller) *MockJWTTokenizer {
	return &MockJWTTokenizer{ctrl: ctrl}
}

func (m *MockJWTTokenizer) Tokenize(claims jwt.Claims) (string, error) {
	return "jwt_token_123", nil
}

func (m *MockJWTTokenizer) GetClaims(token string) (*jwt.Claims, error) {
	return &jwt.Claims{}, nil
}

func TestCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCredentialService := NewMockCredentialService(ctrl)
	mockSessionService := NewMockSessionService(ctrl)
	mockRefreshTokenService := NewMockRefreshTokenService(ctrl)
	mockJWTTokenizer := NewMockJWTTokenizer(ctrl)

	testConfig := &config.Config{
		APPConfig: config.AppConfig{
			CookieDomain: ".weeb.vip",
		},
	}

	t.Run("should create session and set cookies with proper domain", func(t *testing.T) {
		// Set up HTTP response recorder to capture cookies
		recorder := httptest.NewRecorder()
		ctx := responsecontext.WithResponseWriter(context.Background(), recorder)

		input := &model.LoginInput{
			Username: "testuser",
			Password: "testpass",
		}

		// Execute the resolver
		result, err := CreateSession(
			ctx,
			mockCredentialService,
			mockSessionService,
			mockRefreshTokenService,
			mockJWTTokenizer,
			testConfig,
			input,
		)

		// Verify result
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if result.ID != "user_123" {
			t.Errorf("Expected user ID 'user_123', got: %s", result.ID)
		}

		if result.Credentials == nil {
			t.Fatal("Expected credentials, got nil")
		}

		if result.Credentials.Token == nil || *result.Credentials.Token != "jwt_token_123" {
			t.Errorf("Expected token 'jwt_token_123', got: %v", result.Credentials.Token)
		}

		if result.Credentials.RefreshToken == nil || *result.Credentials.RefreshToken != "refresh_token_123" {
			t.Errorf("Expected refresh token 'refresh_token_123', got: %v", result.Credentials.RefreshToken)
		}

		// Verify cookies were set
		cookies := recorder.Result().Cookies()
		if len(cookies) != 2 {
			t.Fatalf("Expected 2 cookies, got: %d", len(cookies))
		}

		// Check access token cookie
		var accessTokenCookie *http.Cookie
		var refreshTokenCookie *http.Cookie

		for _, cookie := range cookies {
			if cookie.Name == "access_token" {
				accessTokenCookie = cookie
			} else if cookie.Name == "refresh_token" {
				refreshTokenCookie = cookie
			}
		}

		if accessTokenCookie == nil {
			t.Fatal("Expected access_token cookie, but not found")
		}

		if accessTokenCookie.Value != "jwt_token_123" {
			t.Errorf("Expected access token cookie value 'jwt_token_123', got: %s", accessTokenCookie.Value)
		}

		if accessTokenCookie.Domain != "weeb.vip" {
			t.Errorf("Expected access token cookie domain 'weeb.vip', got: %s", accessTokenCookie.Domain)
		}

		if !accessTokenCookie.HttpOnly {
			t.Error("Expected access token cookie to be HttpOnly")
		}

		if accessTokenCookie.Path != "/" {
			t.Errorf("Expected access token cookie path '/', got: %s", accessTokenCookie.Path)
		}

		expectedAccessMaxAge := int(time.Hour.Seconds())
		if accessTokenCookie.MaxAge != expectedAccessMaxAge {
			t.Errorf("Expected access token cookie MaxAge %d, got: %d", expectedAccessMaxAge, accessTokenCookie.MaxAge)
		}

		if refreshTokenCookie == nil {
			t.Fatal("Expected refresh_token cookie, but not found")
		}

		if refreshTokenCookie.Value != "refresh_token_123" {
			t.Errorf("Expected refresh token cookie value 'refresh_token_123', got: %s", refreshTokenCookie.Value)
		}

		if refreshTokenCookie.Domain != "weeb.vip" {
			t.Errorf("Expected refresh token cookie domain 'weeb.vip', got: %s", refreshTokenCookie.Domain)
		}

		if !refreshTokenCookie.HttpOnly {
			t.Error("Expected refresh token cookie to be HttpOnly")
		}

		expectedRefreshMaxAge := int((time.Hour * 24 * 7).Seconds())
		if refreshTokenCookie.MaxAge != expectedRefreshMaxAge {
			t.Errorf("Expected refresh token cookie MaxAge %d, got: %d", expectedRefreshMaxAge, refreshTokenCookie.MaxAge)
		}
	})

	t.Run("should use different cookie domain from config", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx := responsecontext.WithResponseWriter(context.Background(), recorder)

		localConfig := &config.Config{
			APPConfig: config.AppConfig{
				CookieDomain: "localhost",
			},
		}

		input := &model.LoginInput{
			Username: "testuser",
			Password: "testpass",
		}

		_, err := CreateSession(
			ctx,
			mockCredentialService,
			mockSessionService,
			mockRefreshTokenService,
			mockJWTTokenizer,
			localConfig,
			input,
		)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		cookies := recorder.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Domain != "localhost" {
				t.Errorf("Expected cookie domain 'localhost', got: %s", cookie.Domain)
			}
		}
	})

	t.Run("should handle guest session when input is nil", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx := responsecontext.WithResponseWriter(context.Background(), recorder)

		// Mock session service to create guest session
		guestSessionService := &MockSessionService{ctrl: ctrl}
		// Override CreateSession for guest scenario

		result, err := CreateSession(
			ctx,
			mockCredentialService,
			guestSessionService,
			mockRefreshTokenService,
			mockJWTTokenizer,
			testConfig,
			nil, // nil input for guest session
		)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// Verify cookies were still set
		cookies := recorder.Result().Cookies()
		if len(cookies) != 2 {
			t.Fatalf("Expected 2 cookies for guest session, got: %d", len(cookies))
		}
	})
}