package credential

const (
	CredentialErrorInternalError      ErrorCode = "INTERNAL_ERROR"      // nolint
	CredentialErrorInvalidCredentials ErrorCode = "INVALID_CREDENTIALS" // nolint
)

type ErrorCode string

type Error struct {
	Code    ErrorCode
	Message string
}

func (c ErrorCode) String() string {
	return string(c)
}

func (e Error) Error() string {
	return e.Message
}
