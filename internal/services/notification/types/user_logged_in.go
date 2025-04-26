// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    userLoggedIn, err := UnmarshalUserLoggedIn(bytes)
//    bytes, err = userLoggedIn.Marshal()

//nolint
package types

import "encoding/json"

func UnmarshalUserLoggedIn(data []byte) (UserLoggedIn, error) {
	var r UserLoggedIn
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UserLoggedIn) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// schema for event which is produced when user logs in.
type UserLoggedIn struct {
	EventName      string `json:"event_name"`       // The type of the event.
	GuestID        string `json:"guest_id"`         // Guest session id before user was logged in.
	LoggedInAt     string `json:"logged_in_at"`     // RFC-3339 formatted time in UTC.
	RefreshTokenID string `json:"refresh_token_id"` // Id of the refresh token user got while logging in.
	SchemaVersion  string `json:"schema_version"`   // The version of the schema for this particular event.
	UserID         string `json:"user_id"`          // The id of the logged in user.
}
