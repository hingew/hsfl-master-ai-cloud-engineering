package auth

type TokenGenerator interface {
	CreateToken(claims map[string]interface{}) (string, error)
	VerifyToken(token string) (map[string]interface{}, error)
}
