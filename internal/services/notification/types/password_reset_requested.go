// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    passwordResetRequested, err := UnmarshalPasswordResetRequested(bytes)
//    bytes, err = passwordResetRequested.Marshal()

//nolint
package types

import "encoding/json"

func UnmarshalPasswordResetRequested(data []byte) (PasswordResetRequested, error) {
	var r PasswordResetRequested
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *PasswordResetRequested) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// A password reset request was created.
type PasswordResetRequested struct {
	EventName      string         `json:"event_name"`      // The type of the event.
	Identifier     string         `json:"identifier"`      // The actual value of user's email or phone number.
	IdentifierType IdentifierType `json:"identifier_type"` // The type of identity user is using for authentication process.
	RequestedAt    string         `json:"requested_at"`    // RFC-3339 formatted time in UTC.
	ResetToken     string         `json:"reset_token"`     // The actual token that can be used by user to verify their email address.
	SchemaVersion  string         `json:"schema_version"`  // The version of the schema for this particular event.
}
