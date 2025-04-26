// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    registrationComplete, err := UnmarshalRegistrationComplete(bytes)
//    bytes, err = registrationComplete.Marshal()

//nolint
package types

import "encoding/json"

func UnmarshalRegistrationComplete(data []byte) (RegistrationComplete, error) {
	var r RegistrationComplete
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *RegistrationComplete) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// A registration was completed.
type RegistrationComplete struct {
	EventName         string            `json:"event_name"`         // The type of the event.
	FirstName         string            `json:"first_name"`         // User's first name.
	Identifier        string            `json:"identifier"`         // The actual value of user's email or phone number.
	IdentifierType    IdentifierType    `json:"identifier_type"`    // The type of identity user is using for authentication process.
	LastName          string            `json:"last_name"`          // User's last name.
	PreferredLanguage PreferredLanguage `json:"preferred_language"` // The language user used while doing the registration.
	RegisteredAt      string            `json:"registered_at"`      // RFC-3339 formatted time in UTC.
	SchemaVersion     string            `json:"schema_version"`     // The version of the schema for this particular event.
	UserID            string            `json:"user_id"`            // User id.
	VerificationToken string            `json:"verification_token"` // The actual token that can be used by user to verify their email address.
}

// the type of identity user is using for authentication process.
type IdentifierType string

const (
	Email IdentifierType = "EMAIL"
)

// the language user used while doing the registration.
type PreferredLanguage string

const (
	En PreferredLanguage = "EN"
	Th PreferredLanguage = "TH"
)
