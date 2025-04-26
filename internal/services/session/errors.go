package session

const (
	SessionErrorInternalError      ErrorCode = "INTERNAL_ERROR"       // nolint
	SessionErrorSessionNotFound    ErrorCode = "SESSION_NOT_FOUND"    // nolint
	SessionErrorInvalidSessionCode ErrorCode = "INVALID_SESSION_CODE" // nolint
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
