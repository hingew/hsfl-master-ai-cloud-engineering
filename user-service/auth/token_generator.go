package auth

type TokenGenerator interface {
	CreateToken(claims map[string]interface{}) (string, error)
}
