// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    refreshTokenExpired, err := UnmarshalRefreshTokenExpired(bytes)
//    bytes, err = refreshTokenExpired.Marshal()

//nolint
package types

import "encoding/json"

func UnmarshalRefreshTokenExpired(data []byte) (RefreshTokenExpired, error) {
	var r RefreshTokenExpired
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RefreshTokenExpired) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// schema for event produced when a refresh token has expired.
type RefreshTokenExpired struct {
	EventName      string `json:"event_name"`       // The type of the event.
	ExpiredAt      string `json:"expired_at"`       // RFC-3339 formatted time in UTC.
	RefreshTokenID string `json:"refresh_token_id"` // Refresh token.
	SchemaVersion  string `json:"schema_version"`   // The version of the schema for this particular event.
	UserID         string `json:"user_id"`          // The user to which refresh token belonged to.
}
