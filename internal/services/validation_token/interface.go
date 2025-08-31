package validation_token // nolint

type ValidationToken interface {
	GenerateToken(identifier string) (string, error)
}
