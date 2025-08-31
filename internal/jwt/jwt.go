package jwt

import (
	"time"

	"github.com/weeb-vip/auth/internal/keypair"

	"github.com/golang-jwt/jwt/v4"
)

const minJWTTokenValidityMinutes = 15

func (t tokenizer) Tokenize(claims Claims) (string, error) {
	signingKey := t.signingKey.GetLatest()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, buildClaims(claims))
	token.Header["kid"] = signingKey.ID

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(signingKey.Key))
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func New(key keypair.RotatingSigningKey) Tokenizer {
	return tokenizer{signingKey: key}
}

func getDefault(duration *time.Duration, defaultDuration time.Duration) time.Duration {
	if duration == nil {
		return defaultDuration
	}

	return *duration
}

func buildClaims(srcClaims Claims) jwt.MapClaims {
	mapClaims := jwt.MapClaims{
		"nbf": time.Now().Unix(),
		"iss": "weeb-vip",
		"aud": "weebusers",
		"iat": time.Now().Unix(),
	}

	mapClaims = addIfNotNil(mapClaims, srcClaims.Subject, "sub")
	mapClaims = addIfNotNil(mapClaims, srcClaims.Purpose, "purpose")
	mapClaims = addIfNotNil(mapClaims, srcClaims.RefreshToken, "refresh_token")
	mapClaims["exp"] = time.
		Now().
		Add(getDefault(srcClaims.TTL, time.Minute*time.Duration(minJWTTokenValidityMinutes))).
		Unix()

	return mapClaims
}

func addIfNotNil[T any](claims jwt.MapClaims, value *T, key string) jwt.MapClaims {
	if value != nil {
		claims[key] = value
	}

	return claims
}

func (t tokenizer) GetClaims(token string) (*Claims, error) {
	signingKey := t.signingKey.GetLatest()
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		publicKey := signingKey

		return jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey.Key))
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return &Claims{
		Subject: mapClaims["sub"].(*string),
		Purpose: mapClaims["purpose"].(*string),
	}, nil
}
