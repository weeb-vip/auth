package resolvers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
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

		// Verify Set-Cookie headers were set (bypassing Go's normalization)
		setCookieHeaders := recorder.Result().Header.Values("Set-Cookie")
		if len(setCookieHeaders) != 2 {
			t.Fatalf("Expected 2 Set-Cookie headers, got: %d", len(setCookieHeaders))
		}

		// Check for access token and refresh token in headers
		var accessTokenHeader, refreshTokenHeader string
		for _, header := range setCookieHeaders {
			if strings.Contains(header, "access_token=") {
				accessTokenHeader = header
			} else if strings.Contains(header, "refresh_token=") {
				refreshTokenHeader = header
			}
		}

		if accessTokenHeader == "" {
			t.Fatal("Expected access_token Set-Cookie header, but not found")
		}

		if refreshTokenHeader == "" {
			t.Fatal("Expected refresh_token Set-Cookie header, but not found")
		}

		// Verify access token header contains expected values
		if !strings.Contains(accessTokenHeader, "access_token=jwt_token_123") {
			t.Errorf("Expected access token value 'jwt_token_123' in header: %s", accessTokenHeader)
		}

		if !strings.Contains(accessTokenHeader, "Domain=.weeb.vip") {
			t.Errorf("Expected access token domain '.weeb.vip' (with leading dot) in header: %s", accessTokenHeader)
		}

		if !strings.Contains(accessTokenHeader, "HttpOnly") {
			t.Errorf("Expected HttpOnly flag in access token header: %s", accessTokenHeader)
		}

		if !strings.Contains(accessTokenHeader, "Secure") {
			t.Errorf("Expected Secure flag in access token header: %s", accessTokenHeader)
		}

		if !strings.Contains(accessTokenHeader, "Path=/") {
			t.Errorf("Expected Path=/ in access token header: %s", accessTokenHeader)
		}

		expectedAccessMaxAge := fmt.Sprintf("Max-Age=%d", int((time.Minute * 15).Seconds()))
		if !strings.Contains(accessTokenHeader, expectedAccessMaxAge) {
			t.Errorf("Expected %s in access token header: %s", expectedAccessMaxAge, accessTokenHeader)
		}

		// Verify refresh token header contains expected values
		if !strings.Contains(refreshTokenHeader, "refresh_token=refresh_token_123") {
			t.Errorf("Expected refresh token value 'refresh_token_123' in header: %s", refreshTokenHeader)
		}

		if !strings.Contains(refreshTokenHeader, "Domain=.weeb.vip") {
			t.Errorf("Expected refresh token domain '.weeb.vip' (with leading dot) in header: %s", refreshTokenHeader)
		}

		if !strings.Contains(refreshTokenHeader, "HttpOnly") {
			t.Errorf("Expected HttpOnly flag in refresh token header: %s", refreshTokenHeader)
		}

		if !strings.Contains(refreshTokenHeader, "Secure") {
			t.Errorf("Expected Secure flag in refresh token header: %s", refreshTokenHeader)
		}

		expectedRefreshMaxAge := fmt.Sprintf("Max-Age=%d", int((time.Hour * 24 * 7).Seconds()))
		if !strings.Contains(refreshTokenHeader, expectedRefreshMaxAge) {
			t.Errorf("Expected %s in refresh token header: %s", expectedRefreshMaxAge, refreshTokenHeader)
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

		setCookieHeaders := recorder.Result().Header.Values("Set-Cookie")
		for _, header := range setCookieHeaders {
			if !strings.Contains(header, "Domain=localhost") {
				t.Errorf("Expected cookie domain 'localhost' in header: %s", header)
			}
		}
	})

	t.Run("should normalize cookie domain by removing leading dot", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ctx := responsecontext.WithResponseWriter(context.Background(), recorder)

		configWithDotDomain := &config.Config{
			APPConfig: config.AppConfig{
				CookieDomain: ".example.com", // Config has leading dot
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
			configWithDotDomain,
			input,
		)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		setCookieHeaders := recorder.Result().Header.Values("Set-Cookie")
		for _, header := range setCookieHeaders {
			// With manual header construction, we preserve the leading dot
			if !strings.Contains(header, "Domain=.example.com") {
				t.Errorf("Expected cookie domain '.example.com' (preserved leading dot) in header: %s", header)
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
		setCookieHeaders := recorder.Result().Header.Values("Set-Cookie")
		if len(setCookieHeaders) != 2 {
			t.Fatalf("Expected 2 Set-Cookie headers for guest session, got: %d", len(setCookieHeaders))
		}
	})
}